// The package man outputs man pages from mmmark markdown.
package man

// Lots of code copied from https://github.com/cpuguy83/go-md2man, but adapated to mmark
// and made to support mmark features.

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/mmarkdown/mmark/v2/lang"
	"github.com/mmarkdown/mmark/v2/mast"
)

// Flags control optional behavior of Markdown renderer.
type Flags int

// HTML renderer configuration options.
const (
	FlagsNone   Flags = 0
	ManFragment Flags = 1 << iota // Don't generate a complete document

	CommonFlags Flags = FlagsNone
)

// RendererOptions is a collection of supplementary parameters tweaking
// the behavior of various parts of Markdown renderer.
type RendererOptions struct {
	Flags Flags // Flags allow customizing this renderer's behavior

	Language lang.Lang // Output language for the document.

	// if set, called at the start of RenderNode(). Allows replacing rendering of some nodes
	RenderNodeHook html.RenderNodeFunc

	// Comments is a list of comments the renderer should detect when
	// parsing code blocks and detecting callouts.
	Comments [][]byte
}

// Renderer implements Renderer interface for Markdown output.
type Renderer struct {
	opts RendererOptions

	Title        *mast.Title
	listLevel    int
	allListLevel int
}

// NewRenderer creates and configures an Renderer object, which satisfies the Renderer interface.
func NewRenderer(opts RendererOptions) *Renderer {
	return &Renderer{opts: opts}
}

func (r *Renderer) hardBreak(w io.Writer, node *ast.Hardbreak) {
	r.outs(w, "\n.br\n")
}

func (r *Renderer) matter(w io.Writer, node *ast.DocumentMatter, entering bool) {
	// TODO: what should this output?
}

func (r *Renderer) title(w io.Writer, node *mast.Title, entering bool) {
	if !entering {
		return
	}

	if node.Date.IsZero() {
		node.Date = time.Now().UTC()
	}

	// track back to first space and assume the rest is the section, don't parse it as a number
	i := len(node.Title) - 1
	for i > 0 && node.Title[i-1] != ' ' {
		i--
	}
	section := 1
	title := node.Title
	switch {
	case i > 0:
		d, err := strconv.Atoi(node.Title[i:])
		if err != nil {
			log.Print("No section number found at end of title, defaulting to 1")
		} else {
			section = d
			title = node.Title[:i-1]
		}
	}
	if i == 0 {
		log.Print("No section number found at end of title, defaulting to 1")
	}

	r.outs(w, fmt.Sprintf(".TH %q", strings.ToUpper(title)))
	r.outs(w, fmt.Sprintf(" %d", section))
	r.outs(w, fmt.Sprintf(" %q", node.Date.Format("January 2006")))
	r.outs(w, fmt.Sprintf(" %q", node.Area))
	r.outs(w, fmt.Sprintf(" %q", node.Workgroup))

	r.outs(w, "\n")
}

func (r *Renderer) heading(w io.Writer, node *ast.Heading, entering bool) {
	if entering {
		switch node.Level {
		case 1, 2:
			r.outs(w, "\n.SH ")
		case 3:
			// normal subsection
			r.outs(w, "\n.SS ")
		default:
			// fake a heading.
			r.outs(w, "\n.PP\n.B ")
		}
	}
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
		// If in lists, suppress paragraphs. Unless we know the list contains
		// block level elements, but then only apply this after the first paragraph.
		parent := para.Parent
		if parent != nil {
			if _, ok := parent.(*ast.ListItem); ok {
				// if we're the first para return, otherwise output a PP
				c := parent.GetChildren()
				i := 0
				par := 0
				for i = range c {
					_, ok := c[i].(*ast.Paragraph)
					if ok {
						par++
					}
					if c[i] == para {
						if par > 1 {
							// No .PP because that messes up formatting.
							r.outs(w, "\n\n")
						}
					}
				}
				return
			}
		}
		r.outs(w, "\n.PP\n")
		return
	}

	r.outs(w, "\n")
}

func (r *Renderer) list(w io.Writer, list *ast.List, entering bool) {
	if list.IsFootnotesList {
		return
	}

	// normal list
	if entering {
		r.allListLevel++
		if list.ListFlags&ast.ListTypeOrdered == 0 && list.ListFlags&ast.ListTypeTerm == 0 && list.ListFlags&ast.ListTypeDefinition == 0 {
			r.listLevel++
		}
		if r.allListLevel > 1 {
			r.outs(w, "\n.RS\n")
		} else {
			r.outs(w, "\n")
		}
		return
	}
	if r.allListLevel > 1 {
		r.outs(w, "\n.RE\n")
	} else {
		r.outs(w, "\n")
	}
	r.allListLevel--
	if list.ListFlags&ast.ListTypeOrdered == 0 && list.ListFlags&ast.ListTypeTerm == 0 && list.ListFlags&ast.ListTypeDefinition == 0 {
		r.listLevel--
	}
}

func (r *Renderer) listItem(w io.Writer, listItem *ast.ListItem, entering bool) {
	if entering {
		// footnotes
		if listItem.RefLink != nil {
			// get number in the list
			children := listItem.Parent.GetChildren()
			for i := range children {
				if listItem == children[i] {
					r.outs(w, fmt.Sprintf("\n.IP [%d]\n", i+1))
				}
			}
			return
		}

		x := listItem.ListFlags
		switch {
		case x&ast.ListTypeOrdered != 0:
			children := listItem.GetParent().GetChildren()
			i := 0
			for i = 0; i < len(children); i++ {
				if children[i] == listItem {
					break
				}
			}
			start := listItem.GetParent().(*ast.List).Start
			r.outs(w, fmt.Sprintf(".IP %d\\. 4\n", start+i+1))

		case x&ast.ListTypeTerm != 0:
			r.outs(w, ".TP\n")

		case x&ast.ListTypeDefinition != 0:
			r.outs(w, "")

		default:
			if r.listLevel%2 == 0 {
				r.outs(w, ".IP \\(en 4\n")
			} else {
				r.outs(w, ".IP \\(bu 4\n")
			}
		}
	}
}

func (r *Renderer) codeBlock(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	if entering {
		r.outs(w, "\n.PP\n.RS\n\n.nf\n")
		escapeSpecialChars(r, w, codeBlock.Literal)
		r.outs(w, "\n.fi\n.RE\n")
	}
}

func (r *Renderer) table(w io.Writer, tab *ast.Table, entering bool) {
	// The tbl renderer want to see the entire table's columns, rows first
	if entering {
		r.outs(w, "\n.RS\n.TS\nallbox;\n")
		cells := rows(tab)
		for r1 := 0; r1 < len(cells); r1++ {
			align := ""
			for c := 0; c < len(cells[r1]); c++ {
				x := cells[r1][c]
				switch x.Align {
				case ast.TableAlignmentLeft:
					align += "l "
				case ast.TableAlignmentRight:
					align += "r "
				case ast.TableAlignmentCenter:
					fallthrough
				default:
					align += "c "
				}
				if x.ColSpan > 0 {
					align += strings.Repeat("s ", x.ColSpan-1)
				}

			}
			r.outs(w, strings.TrimSpace(align)+"\n")
		}
		r.outs(w, ".\n")
		return
	}
	r.outs(w, ".TE\n.RE\n\n")
}

func (r *Renderer) tableRow(w io.Writer, tableRow *ast.TableRow, entering bool) {
	if !entering {
		r.outs(w, "\n")
	}
}

func (r *Renderer) tableCell(w io.Writer, tableCell *ast.TableCell, entering bool) {
	if tableCell.IsHeader {
		r.outOneOf(w, entering, "\\fB", "\\fP")
	}
	parent := tableCell.Parent
	if tableCell == ast.GetFirstChild(parent) {
		return
	}
	if entering {
		r.outs(w, "\t")
		return
	}
}

func (r *Renderer) htmlSpan(w io.Writer, span *ast.HTMLSpan) {}

func (r *Renderer) crossReference(w io.Writer, cr *ast.CrossReference, entering bool) {
	if !entering {
		return
	}
	r.out(w, bytes.ToUpper(cr.Destination))
}

func (r *Renderer) index(w io.Writer, index *ast.Index, entering bool) {}

func (r *Renderer) link(w io.Writer, link *ast.Link, entering bool) {
	if link.Footnote != nil {
		if entering {
			r.outs(w, fmt.Sprintf("\\u[%d]\\d", link.NoteID))
		}
		return
	}
	// !entering so the URL comes after the link text.
	if !entering {
		r.outs(w, "\n\\[la]")
		r.out(w, link.Destination)
		r.outs(w, "\\[ra]")
	}
}

func (r *Renderer) image(w io.Writer, node *ast.Image, entering bool) {
	// only works with `ascii-art` images
	if !bytes.HasSuffix(node.Destination, []byte(".ascii-art")) {
		// remove from tree, we can use RemoveFromTree, because that re-shuffles things and can
		// make us output things twice.
		node.SetChildren(nil)
		node = nil
		return
	}
	if !entering {
		r.outs(w, "\n.fi\n.RE\n")
		return
	}
	node.SetChildren(nil) // remove Title, if any, we can type set it.
	r.outs(w, "\n.PP\n.RS\n\n.nf\n")
	img, err := ioutil.ReadFile(string(node.Destination)) // markdown, doens't err, this can be an empty image, log maybe??
	if err != nil {
		img = []byte(err.Error())
	}
	escapeSpecialChars(r, w, img)
}

func (r *Renderer) mathBlock(w io.Writer, mathBlock *ast.MathBlock, entering bool) {
	// may indent it?
}

func (r *Renderer) captionFigure(w io.Writer, figure *ast.CaptionFigure, entering bool) {
	// check what we have here and throw away any non ascii-art figures
	// CaptionFigure
	//   Paragraph
	//     Text
	//     Image 'url=array-vs-slice.svg'
	//       Text 'Array versus Slice svg'
	//     Text '\n'
	//     Image 'url=array-vs-slice.ascii-art'
	//       Text 'Array versus Slice ascii-art'
	//
	// The image with svg will be removed as child and then we continue to walk the AST.
	for _, child := range figure.GetChildren() {
		// for figures/images, these are wrapped below a paragraph.
		// TODO: can there be more than 1 paragraph??
		if p, ok := child.(*ast.Paragraph); ok {
			for _, img := range p.GetChildren() {
				x, ok := img.(*ast.Image)
				if !ok {
					continue
				}
				// if not ascii-art, remove entirely
				if !bytes.HasSuffix(x.Destination, []byte(".ascii-art")) {
					ast.RemoveFromTree(img) // this is save because we're not accessing any of the children just yet.
					continue
				}
				img.SetChildren(nil) // remove alt text
			}
		}
	}
}

func (r *Renderer) caption(w io.Writer, caption *ast.Caption, entering bool) {
	what := ast.GetFirstChild(caption.Parent)

	if !entering {
		switch what.(type) {
		case *ast.Table:
			r.outs(w, "\n.RE\n")
		case *ast.CodeBlock, *ast.Paragraph: // Paragraph is here because that wrap an image.
			r.outs(w, "\n.RE\n")
		case *ast.BlockQuote:
			r.outs(w, "\n.RE\n")
		}
		return
	}
	// get parent, get first child for type
	switch what.(type) {
	case *ast.Table:
		r.outs(w, "\n.RS\n")
	case *ast.CodeBlock, *ast.Paragraph:
		r.outs(w, "\n.RS\n")
	case *ast.BlockQuote:
		r.outs(w, "\n.RS\n")
		r.outs(w, "\\(en ")
	}
}

func (r *Renderer) blockQuote(w io.Writer, block *ast.BlockQuote, entering bool) {
	if entering {
		r.outs(w, "\n.PP\n.RS\n")
	} else {
		r.outs(w, "\n.RE\n")
	}
}

func (r *Renderer) aside(w io.Writer, block *ast.Aside, entering bool) {
	if entering {
		r.outs(w, "\n.PP\n.RS\n")
	} else {
		r.outs(w, "\n.RE\n")
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

	if attr := mast.AttributeFromNode(node); attr != nil && entering {
	}

	switch node := node.(type) {
	case *ast.Document:
		// do nothing
	case *mast.Title:
		r.title(w, node, entering)
		r.Title = node // save for later.
	case *mast.Authors:
		r.authors(w, node, entering)
	case *mast.Bibliography, *mast.BibliographyWrapper:
		if entering {
			r.outs(w, "\n.SH \"")
			r.outs(w, strings.ToUpper(r.opts.Language.Bibliography()))
			r.outs(w, "\"\n")
		}
	case *mast.BibliographyItem:
		r.bibliographyItem(w, node, entering)
	case *mast.DocumentIndex, *mast.IndexLetter, *mast.IndexItem, *mast.IndexSubItem, *mast.IndexLink:
	case *mast.ReferenceBlock:
		// ignore
	case *ast.Footnotes:
		r.footnotes(w, node, entering)
	case *ast.Text:
		r.text(w, node, entering)
	case *ast.Softbreak:
		// TODO
	case *ast.Hardbreak:
		r.hardBreak(w, node)
	case *ast.NonBlockingSpace:
		r.outs(w, "\\0")
	case *ast.Callout:
		r.callout(w, node, entering)
	case *ast.Emph:
		r.outOneOf(w, entering, "\\fI", "\\fP")
	case *ast.Strong:
		r.outOneOf(w, entering, "\\fB", "\\fP")
	case *ast.Del:
		r.outOneOf(w, entering, "~~", "~~")
	case *ast.Citation:
		r.citation(w, node, entering)
	case *ast.DocumentMatter:
		r.matter(w, node, entering)
	case *ast.Heading:
		r.heading(w, node, entering)
	case *ast.HorizontalRule:
		if entering {
			r.outs(w, "\n.ti 0\n\\l'\\n(.l─'\n")
		}
	case *ast.Paragraph:
		r.paragraph(w, node, entering)
	case *ast.HTMLSpan:
		r.out(w, node.Literal)
	case *ast.HTMLBlock:
		r.out(w, node.Literal)
	case *ast.List:
		r.list(w, node, entering)
	case *ast.ListItem:
		r.listItem(w, node, entering)
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
	case *ast.TableBody:
	case *ast.TableFooter:
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
		if entering {
			r.out(w, node.Literal)
		}
	case *ast.Image:
		r.image(w, node, entering)
	case *ast.Code:
		r.outs(w, "\\fB\\fC")
		r.out(w, node.Literal)
		r.outs(w, "\\fR")
	case *ast.MathBlock:
		r.mathBlock(w, node, entering)
	case *ast.Subscript:
		r.outOneOf(w, true, "\\d", "\\u")
		if entering {
			r.out(w, node.Literal)
		}
		r.outOneOf(w, false, "\\d", "\\u")
	case *ast.Superscript:
		r.outOneOf(w, true, "\\u", "\\d")
		if entering {
			r.out(w, node.Literal)
		}
		r.outOneOf(w, false, "\\u", "\\d")
	default:
		panic(fmt.Sprintf("Unknown node %T", node))
	}
	return ast.GoToNext
}

func (r *Renderer) callout(w io.Writer, node *ast.Callout, entering bool) {
	if entering {
		r.outs(w, "\\fB")
		r.out(w, node.ID)
		r.outs(w, "\\fP")
		return
	}
}

func (r *Renderer) text(w io.Writer, node *ast.Text, entering bool) {
	if !entering {
		return
	}
	text := node.Literal
	parent := node.Parent
	if parent != nil {
		if _, ok := parent.(*ast.Heading); ok {
			text = bytes.ToUpper(text)
			text = append(text, byte('"'))
			text = append([]byte{byte('"')}, text...)
		}
	}

	// A hyphen must be escaped otherwise it will be an em/en-dash.
	text = bytes.Replace(text, []byte("-"), []byte("\\-"), -1)

	r.out(w, text)
}

func (r *Renderer) footnotes(w io.Writer, node ast.Node, entering bool) {
	if !entering {
		return
	}
	r.outs(w, "\n.SH \""+strings.ToUpper(r.opts.Language.Footnotes())+"\"\n")
}

func (r *Renderer) RenderHeader(w io.Writer, _ ast.Node) {
	if r.opts.Flags&ManFragment != 0 {
		return
	}
	r.outs(w, `.\" Generated by Mmark Markdown Processer - mmark.miek.nl`+"\n")
}

func (r *Renderer) RenderFooter(w io.Writer, node ast.Node) {}

func (r *Renderer) bibliographyItem(w io.Writer, bib *mast.BibliographyItem, entering bool) {
	if !entering {
		return
	}
	if bib.Reference == nil {
		return
	}
	r.outs(w, ".TP\n")
	r.outs(w, fmt.Sprintf("[%s]\n", bib.Anchor))
	for _, author := range bib.Reference.Front.Authors {
		writeNonEmptyString(w, author.Fullname)
		if author.Organization != nil {
			writeNonEmptyString(w, author.Organization.Value)
		}
	}

	writeNonEmptyString(w, bib.Reference.Front.Title)
	if bib.Reference.Target != "" {
		r.outs(w, "\\[la]")
		r.outs(w, bib.Reference.Target)
		r.outs(w, "\\[ra]")
	}
	writeNonEmptyString(w, bib.Reference.Front.Date.Year)
	r.outs(w, "\n")
}

func writeNonEmptyString(w io.Writer, s string) {
	if s == "" {
		return
	}
	io.WriteString(w, s)
	io.WriteString(w, "\n")
}
