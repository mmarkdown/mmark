// The package text outputs text.
package text

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/mmarkdown/mmark/v2/mast"
)

// Flags control optional behavior of Markdown renderer.
type Flags int

// Markdown renderer configuration options.
const (
	FlagsNone Flags = 0

	CommonFlags Flags = FlagsNone
)

// RendererOptions is a collection of supplementary parameters tweaking
// the behavior of various parts of Markdown renderer.
type RendererOptions struct {
	Flags Flags // Flags allow customizing this renderer's behavior

	TextWidth int

	// if set, called at the start of RenderNode(). Allows replacing rendering of some nodes
	RenderNodeHook html.RenderNodeFunc
}

// Renderer implements Renderer interface for Markdown output.
type Renderer struct {
	opts RendererOptions

	headingTransformFunc func([]byte) []byte // How do display heading, noop if nothing needs doing

	// TODO(miek): paraStart should probably be a stack, aside in para in aside, etc.
	paraStart    int
	headingStart int

	prefix *prefixStack // track current prefix, quote, aside, etc.

	// tables
	cellStart int
	col       int
	colWidth  []int
	colAlign  []ast.CellAlignFlags
	tableType ast.Node

	suppress bool // when true we suppress newlines

	deferredFootBuf *bytes.Buffer // deferred footnote buffer. Appended to the doc at the end.
	deferredFootID  map[string]struct{}

	deferredLinkBuf *bytes.Buffer // deferred footnote buffer. Appended to the doc at the end.
	deferredLinkID  map[string]struct{}
}

// NewRenderer creates and configures an Renderer object, which satisfies the Renderer interface.
func NewRenderer(opts RendererOptions) *Renderer {
	if opts.TextWidth == 0 {
		opts.TextWidth = 80
	}
	r := &Renderer{
		opts:                 opts,
		prefix:               &prefixStack{p: [][]byte{}},
		deferredFootBuf:      &bytes.Buffer{},
		deferredFootID:       make(map[string]struct{}),
		deferredLinkBuf:      &bytes.Buffer{},
		deferredLinkID:       make(map[string]struct{}),
		headingTransformFunc: noopHeadingTransferFunc,
	}
	r.push(Space(0)) // default indent for all text, except heading.
	return r
}

func (r *Renderer) hardBreak(w io.Writer, node *ast.Hardbreak) {
	r.endline(w)
	r.newline(w)
}

func (r *Renderer) heading(w io.Writer, node *ast.Heading, entering bool) {
	if entering {
		switch node.Level {
		case 1:
			r.headingTransformFunc = func(data []byte) []byte {
				x := r.centerText(bytes.ToUpper(data))
				return x
			}
		case 2:
			r.headingTransformFunc = func(data []byte) []byte {
				x := r.centerText(data)
				return x
			}
		default:
			r.headingTransformFunc = noopHeadingTransferFunc
		}
		return
	}

	r.headingTransformFunc = noopHeadingTransferFunc
	r.newline(w)
	r.newline(w)
	if node.Level == 1 {
		r.newline(w)
	}
	return
}

func (r *Renderer) horizontalRule(w io.Writer, node *ast.HorizontalRule) {
	r.newline(w)
	r.outs(w, "******")
	r.newline(w)
}

func (r *Renderer) citation(w io.Writer, node *ast.Citation, entering bool) {
	r.outs(w, "[")
	for i, dest := range node.Destination {
		if i > 0 {
			r.outs(w, ", ")
		}
		r.out(w, dest)

	}
	r.outs(w, "]")
}

func (r *Renderer) paragraph(w io.Writer, para *ast.Paragraph, entering bool) {
	if entering {
		if buf, ok := w.(*bytes.Buffer); ok {
			r.paraStart = buf.Len()
		}
		return
	}

	buf, ok := w.(*bytes.Buffer)
	end := 0
	if ok {
		end = buf.Len()
	}
	// Reformat the entire buffer and rewrite to the writer.
	b := buf.Bytes()[r.paraStart:end]

	var indented []byte
	p := bytes.Split(b, []byte("\\\n"))
	for i := range p {
		if len(indented) > 0 {
			p1 := r.wrapText(p[i], r.prefix.flatten(), []byte(""))
			indented = append(indented, []byte("\\\n")...)
			indented = append(indented, p1...)
			continue
		}
		indented = r.wrapText(p[i], r.prefix.flatten(), []byte(""))
	}
	if len(indented) == 0 {
		indented = make([]byte, r.prefix.peek()+3)
	}

	buf.Truncate(r.paraStart)

	// Now an indented list didn't get is marker yet, override the initial spaces that have been
	// created with the list marker, taking the current prefix into account.
	listItem, inList := para.Parent.(*ast.ListItem)
	firstPara := ast.GetPrevNode(para) // only the first para in the listItem needs a list marker
	if inList && firstPara == nil {
		plen := r.prefix.len() - r.prefix.peek()
		switch x := listItem.ListFlags; {
		case x&ast.ListTypeOrdered != 0:
			list := listItem.Parent.(*ast.List) // this must be always true
			pos := []byte(strconv.Itoa(list.Start))
			for i := 0; i < len(pos); i++ {
				indented[plen+i] = pos[i]
			}
			indented[plen+len(pos)] = '.'
			indented[plen+len(pos)+1] = ' '

			list.Start++
		case x&ast.ListTypeTerm != 0:
			indented = indented[plen+r.prefix.peek()-3:] // remove prefix.
			indented[plen+0] = '*'
		case x&ast.ListTypeDefinition != 0:
			indented[plen+0] = ' '
			indented[plen+1] = ' '
			indented[plen+2] = ' '
		default:
			if plen == 0 {
				indented[plen+0] = 'o'
			}
			if plen == 3 {
				indented[plen+0] = '+'
			}
			if plen == 6 {
				indented[plen+0] = 'o'
			}
			if plen > 6 {
				indented[plen+0] = '-'
			}
			indented[plen+1] = ' '
			indented[plen+2] = ' '
		}
	}

	r.out(w, indented)
	r.endline(w)

	// A paragraph can be rendered if we are in a subfigure, if so suppress some newlines.
	if _, inCaption := para.Parent.(*ast.CaptionFigure); inCaption {
		return
	}
	if !lastNode(para) {
		r.newline(w)
	}
}

func (r *Renderer) list(w io.Writer, list *ast.List, entering bool) {
	if entering {
		if list.Start == 0 {
			list.Start = 1
		}
		l := listPrefixLength(list, list.Start)
		if list.ListFlags&ast.ListTypeOrdered != 0 {
			r.push(Space(l))
		} else {
			r.push(Space(3))
		}
		return
	}
	r.pop()
}

func (r *Renderer) codeBlock(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	// Indent codeblock with 3 spaces
	indented := r.indentText(codeBlock.Literal, append(r.prefix.flatten(), Space(3)...))
	r.out(w, indented)
	r.outPrefix(w)
	r.newline(w)

	r.newline(w)
	return
}

func (r *Renderer) table(w io.Writer, tab *ast.Table, entering bool) {
	if entering {
		r.colWidth, r.colAlign = r.tableColWidth(tab)
		r.col = 0
	} else {
		r.colWidth = []int{}
		r.colAlign = []ast.CellAlignFlags{}
	}
}

func (r *Renderer) tableRow(w io.Writer, tableRow *ast.TableRow, entering bool) {
	if entering {
		r.outPrefix(w)
		r.col = 0
		for i, width := range r.colWidth {
			if _, isFooter := r.tableType.(*ast.TableFooter); isFooter {
				r.out(w, bytes.Repeat([]byte("="), width+1))

				if i == len(r.colWidth)-1 {
					r.endline(w)
					r.outPrefix(w)
				} else {
					r.outs(w, "|")
				}
			}
		}

		return
	}

	for i, width := range r.colWidth {
		if _, isHeader := r.tableType.(*ast.TableHeader); isHeader {
			if i == 0 {
				r.outPrefix(w)
			}
			heading := bytes.Repeat([]byte("-"), width+1)

			switch r.colAlign[i] {
			case ast.TableAlignmentLeft:
				heading[0] = '-'
			case ast.TableAlignmentRight:
				heading[width] = '-'
			}
			r.out(w, heading)
			if i == len(r.colWidth)-1 {
				r.endline(w)
			} else {
				r.outs(w, "|")
			}
		}
	}
}

func (r *Renderer) tableCell(w io.Writer, tableCell *ast.TableCell, entering bool) {
	// we get called when we're calculating the column width, only when r.tableColWidth is set we need to output.
	if len(r.colWidth) == 0 {
		return
	}
	if entering {
		if buf, ok := w.(*bytes.Buffer); ok {
			r.cellStart = buf.Len() + 1
		}
		if r.col > 0 {
			r.out(w, Space1)
		}
		return
	}

	cur := 0
	if buf, ok := w.(*bytes.Buffer); ok {
		cur = buf.Len()
	}
	size := r.colWidth[r.col]
	fill := bytes.Repeat(Space1, size-(cur-r.cellStart))
	r.out(w, fill)
	if r.col == len(r.colWidth)-1 {
		r.endline(w)
	} else {
		r.outs(w, "|")
	}
	r.col++
}

func (r *Renderer) htmlSpan(w io.Writer, span *ast.HTMLSpan) {}

func (r *Renderer) crossReference(w io.Writer, cr *ast.CrossReference, entering bool) {
	if entering {
		r.outs(w, "(#")
		r.out(w, cr.Destination)
		return
	}
	r.outs(w, ")")
}

func (r *Renderer) index(w io.Writer, index *ast.Index, entering bool) {
	if !entering {
		return
	}

	r.outs(w, "(!")
	if index.Primary {
		r.outs(w, "!")
	}
	r.out(w, index.Item)

	if len(index.Subitem) > 0 {
		r.outs(w, ", ")
		r.out(w, index.Subitem)
	}
	r.outs(w, ")")
}

func (r *Renderer) link(w io.Writer, link *ast.Link, entering bool) {
	if !entering {
		return
	}
	// clear link so we don't render any children.
	defer func() { *link = ast.Link{} }()

	// footnote
	if link.NoteID > 0 {
		ast.RemoveFromTree(link.Footnote)
		if len(link.DeferredID) > 0 {

			r.outs(w, "[^")
			r.out(w, link.DeferredID)
			r.outs(w, "]")

			if _, ok := r.deferredFootID[string(link.DeferredID)]; ok {
				return
			}

			r.deferredFootBuf.Write(Space(3))
			r.deferredFootBuf.Write([]byte("[^"))
			r.deferredFootBuf.Write(link.DeferredID)
			r.deferredFootBuf.Write([]byte("]: "))
			r.deferredFootBuf.Write(link.Title)

			r.deferredFootID[string(link.DeferredID)] = struct{}{}

			return
		}
		r.outs(w, "^[")
		r.out(w, link.Title)
		r.outs(w, "]")
		return
	}

	for _, child := range link.GetChildren() {
		ast.WalkFunc(child, func(node ast.Node, entering bool) ast.WalkStatus {
			if text, ok := node.(*ast.Text); ok {
				if bytes.Compare(text.Literal, link.Destination) == 0 {
					return ast.GoToNext
				}
			}
			return r.RenderNode(w, node, entering)
		})
	}

	if len(link.DeferredID) == 0 {
		r.outs(w, " <")
		r.out(w, link.Destination)
		r.outs(w, ">")
		link.Destination = []byte{}

		if len(link.Title) > 0 {
			r.outs(w, ` "`)
			r.out(w, link.Title)
			r.outs(w, `"`)
		}
		return
	}

	r.outs(w, "[")
	r.out(w, link.DeferredID)
	r.outs(w, "]")

	if _, ok := r.deferredLinkID[string(link.DeferredID)]; ok {
		return
	}

	r.out(r.deferredLinkBuf, Space(3))
	r.outs(r.deferredLinkBuf, "[")
	r.out(r.deferredLinkBuf, link.DeferredID)
	r.outs(r.deferredLinkBuf, "]: ")
	r.out(r.deferredLinkBuf, link.Destination)
	if len(link.Title) > 0 {
		r.outs(r.deferredLinkBuf, ` "`)
		r.out(r.deferredLinkBuf, link.Title)
		r.outs(r.deferredLinkBuf, `"`)
	}

	r.deferredLinkID[string(link.DeferredID)] = struct{}{}
}

func (r *Renderer) image(w io.Writer, node *ast.Image, entering bool) {
	if !entering {
		return
	}
	// clear image so we don't render any children.
	defer func() { *node = ast.Image{} }()

	r.outs(w, "![")
	for _, child := range node.GetChildren() {
		ast.WalkFunc(child, func(node ast.Node, entering bool) ast.WalkStatus {
			return r.RenderNode(w, node, entering)
		})
	}
	r.outs(w, "]")

	r.outs(w, "(")
	r.out(w, node.Destination)
	if len(node.Title) > 0 {
		r.outs(w, ` "`)
		r.out(w, node.Title)
		r.outs(w, `"`)
	}
	r.outs(w, ")")
}

func (r *Renderer) mathBlock(w io.Writer, mathBlock *ast.MathBlock, entering bool) {
	if !entering {
		return
	}
	r.outPrefix(w)
	r.outs(w, "$$")

	math := r.indentText(mathBlock.Literal, r.prefix.flatten())
	r.out(w, math)

	r.outPrefix(w)
	r.outs(w, "$$\n")

	if !lastNode(mathBlock) {
		r.newline(w)
	}
}

func (r *Renderer) captionFigure(w io.Writer, figure *ast.CaptionFigure, entering bool) {
	// if one of our children is an image this is an subfigure.
	isImage := false
	ast.WalkFunc(figure, func(node ast.Node, entering bool) ast.WalkStatus {
		_, isImage = node.(*ast.Image)
		if isImage {
			return ast.Terminate
		}
		return ast.GoToNext

	})
	if isImage && entering {
		r.outs(w, "")
		r.endline(w)
	}
	if !entering {
		r.newline(w)
	}
}

func (r *Renderer) caption(w io.Writer, caption *ast.Caption, entering bool) {
	if !entering {
		r.endline(w)
		r.newline(w)
		return
	}

	r.outPrefix(w)
	switch ast.GetPrevNode(caption).(type) {
	case *ast.BlockQuote:
		r.outs(w, "Quote: ")
		return
	case *ast.Table:
		r.outs(w, "Table: ")
		return
	case *ast.CodeBlock:
		r.outs(w, "Figure: ")
		return
	}
	// If here, we're dealing with a subfigure captionFigure.
	r.outs(w, "")
	r.endline(w)
	r.outs(w, "")
}

func (r *Renderer) blockQuote(w io.Writer, block *ast.BlockQuote, entering bool) {
	if entering {
		r.push(Quote)
		return
	}
	r.pop()
	r.newline(w)
}

func (r *Renderer) aside(w io.Writer, block *ast.Aside, entering bool) {
	if entering {
		r.push(Aside)
		return
	}
	r.pop()
	if !lastNode(block) {
		r.newline(w)
	}
}

// RenderNode renders a markdown node to markdown.
func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	if r.opts.RenderNodeHook != nil {
		status, didHandle := r.opts.RenderNodeHook(w, node, entering)
		if didHandle {
			return status
		}
	}

	switch node := node.(type) {
	case *ast.Document:
		// do nothing
	case *mast.Title:
		r.out(w, node.Content)
		r.outs(w, "\n")
		r.newline(w)
	case *mast.Bibliography:
		// discard
	case *mast.BibliographyItem:
		// discard
	case *mast.DocumentIndex, *mast.IndexLetter, *mast.IndexItem, *mast.IndexSubItem, *mast.IndexLink:
		// discard
	case *mast.ReferenceBlock:
		// discard
	case *ast.Footnotes:
		// do nothing, we're not outputing a footnote list
	case *ast.Text:
		r.text(w, node, entering)
	case *ast.Softbreak:
	case *ast.Hardbreak:
		r.hardBreak(w, node)
	case *ast.Callout:
		r.callout(w, node, entering)
	case *ast.Emph:
		r.out(w, node.Literal)
	case *ast.Strong:
		r.out(w, node.Literal)
	case *ast.Del:
		r.out(w, node.Literal)
	case *ast.Citation:
		r.citation(w, node, entering)
	case *ast.DocumentMatter:
		// don't output
	case *ast.Heading:
		r.heading(w, node, entering)
	case *ast.HorizontalRule:
		if entering {
			r.newline(w)
			r.outPrefix(w)
			r.outs(w, "********\n")
			r.newline(w)
		}
	case *ast.Paragraph:
		r.paragraph(w, node, entering)
	case *ast.HTMLSpan:
		r.out(w, node.Literal)
	case *ast.HTMLBlock:
		r.out(w, node.Literal)
		r.endline(w)
		r.newline(w)
	case *ast.List:
		r.list(w, node, entering)
		if !entering {
			r.newline(w)
		}
	case *ast.ListItem:
	case *ast.CodeBlock:
		r.codeBlock(w, node, entering)
	case *ast.Caption:
		r.caption(w, node, entering)
	case *ast.CaptionFigure:
		r.captionFigure(w, node, entering)
	case *ast.Table:
		r.table(w, node, entering)
	case *ast.TableCell:
		r.tableCell(w, node, entering)
	case *ast.TableHeader:
		r.tableType = node
	case *ast.TableBody:
		r.tableType = node
	case *ast.TableFooter:
		r.tableType = node
	case *ast.TableRow:
		r.tableRow(w, node, entering)
	case *ast.BlockQuote:
		r.blockQuote(w, node, entering)
	case *ast.Aside:
		r.aside(w, node, entering)
	case *ast.CrossReference:
		r.crossReference(w, node, entering)
	case *ast.Index:
		r.index(w, node, entering)
	case *ast.Link:
		r.link(w, node, entering)
	case *ast.Math:
		r.outOneOf(w, true, "$", "$")
		if entering {
			r.out(w, node.Literal)
		}
		r.outOneOf(w, false, "$", "$")
	case *ast.Image:
		r.image(w, node, entering)
	case *ast.Code:
		r.outs(w, "`")
		r.out(w, node.Literal)
		r.outs(w, "`")
	case *ast.MathBlock:
		r.mathBlock(w, node, entering)
	case *ast.Subscript:
		r.outOneOf(w, true, "~", "~")
		if entering {
			r.out(w, node.Literal)
		}
		r.outOneOf(w, false, "~", "~")
	case *ast.Superscript:
		r.outOneOf(w, true, "^", "^")
		if entering {
			r.out(w, node.Literal)
		}
		r.outOneOf(w, false, "^", "^")
	default:
		panic(fmt.Sprintf("Unknown node %T", node))
	}
	return ast.GoToNext
}

func (r *Renderer) callout(w io.Writer, node *ast.Callout, entering bool) {
	if !entering {
		return
	}
	r.outs(w, "<<")
	r.out(w, node.ID)
	r.outs(w, ">>")
}

func (r *Renderer) text(w io.Writer, node *ast.Text, entering bool) {
	if !entering {
		return
	}
	_, isTableCell := node.Parent.(*ast.TableCell)
	if isTableCell {
		allSpace := true
		for i := range node.Literal {
			if !isSpace(node.Literal[i]) {
				allSpace = false
				break
			}
		}
		if allSpace {
			return
		}
	}

	r.out(w, node.Literal)
}

func (r *Renderer) RenderHeader(_ io.Writer, _ ast.Node) {}
func (r *Renderer) writeDocumentHeader(_ io.Writer)      {}

func (r *Renderer) RenderFooter(w io.Writer, _ ast.Node) {
	if r.deferredFootBuf.Len() > 0 {
		r.outs(w, "\n")
		io.Copy(w, r.deferredFootBuf)
	}
	if r.deferredLinkBuf.Len() > 0 {
		r.outs(w, "\n")
		io.Copy(w, r.deferredLinkBuf)
	}

	buf, ok := w.(*bytes.Buffer)
	if !ok {
		return
	}

	trimmed := &bytes.Buffer{}

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		trimmed.Write(bytes.TrimRight(scanner.Bytes(), " "))
		trimmed.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return
	}

	buf.Truncate(0)
	data := trimmed.Bytes()
	ld := len(data)
	if ld > 2 && data[ld-1] == '\n' && data[ld-2] == '\n' {
		ld--
	}
	buf.Write(data[:ld])
}

var (
	Space1 = Space(1)
	Aside  = []byte("| ")
	Quote  = []byte("| ")
)
