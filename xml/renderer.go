package xml

import (
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/mmarkdown/mmark/mast"
)

// Flags control optional behavior of XML3 renderer.
type Flags int

// HTML renderer configuration options.
const (
	FlagsNone   Flags = iota
	XMLFragment       // Don't generate a complete XML document
	SkipHTML          // Skip preformatted HTML blocks - skips comments

	CommonFlags Flags = FlagsNone
)

type RendererOptions struct {
	// Callouts are supported and detected by setting this option to the callout prefix.
	Callout string

	Flags Flags // Flags allow customizing this renderer's behavior

	// if set, called at the start of RenderNode(). Allows replacing
	// rendering of some nodes
	RenderNodeHook html.RenderNodeFunc

	// Comments is a list of comments the renderer should detect when
	// parsing code blocks and detecting callouts.
	Comments [][]byte
}

// Renderer implements Renderer interface for IETF XMLv3 output. See RFC 7991.
type Renderer struct {
	opts RendererOptions

	documentMatter ast.DocumentMatters // keep track of front/main/back matter
	section        *ast.Heading        // current open section
}

// New creates and configures an Renderer object, which satisfies the Renderer interface.
func NewRenderer(opts RendererOptions) *Renderer {
	return &Renderer{opts: opts}
}

func (r *Renderer) text(w io.Writer, text *ast.Text) {
	if _, parentIsLink := text.Parent.(*ast.Link); parentIsLink {
		//html.EscLink(w, text.Literal)
		r.out(w, text.Literal)
		return
	}
	if heading, parentIsHeading := text.Parent.(*ast.Heading); parentIsHeading {
		if isAbstract(heading.Special) {
			// No <name> when abstract, should not output anything
			return
		}
		r.outs(w, "<name>")
		html.EscapeHTML(w, text.Literal)
		r.outs(w, "</name>")
		return
	}
	if _, parentIsReferences := text.Parent.(*mast.References); parentIsReferences {
		r.outs(w, "<name>")
		html.EscapeHTML(w, text.Literal)
		r.outs(w, "</name>")
		return
	}

	html.EscapeHTML(w, text.Literal)
}

func (r *Renderer) hardBreak(w io.Writer, node *ast.Hardbreak) {
	r.outs(w, "<br />")
	r.cr(w)
}

func (r *Renderer) strong(w io.Writer, node *ast.Strong, entering bool) {
	// *iff* we have a text node as a child *and* that text is 2119, we output bcp14 tags, otherwise just string.
	text := ast.GetFirstChild(node)
	if t, ok := text.(*ast.Text); ok {
		if is2119(t.Literal) {
			r.outOneOf(w, entering, "<bcp14>", "</bcp14>")
			return
		}
	}

	r.outOneOf(w, entering, "<strong>", "</strong>")
}

func (r *Renderer) matter(w io.Writer, node *ast.DocumentMatter) {
	r.sectionClose(w)
	r.section = nil

	switch node.Matter {
	case ast.DocumentMatterFront:
		r.cr(w)
		r.outs(w, "<front>")
		r.cr(w)
	case ast.DocumentMatterMain:
		r.cr(w)
		r.outs(w, "</front>")
		r.cr(w)
		r.cr(w)
		r.outs(w, "<middle>")
		r.cr(w)
	case ast.DocumentMatterBack:
		r.cr(w)
		r.outs(w, "</middle>")
		r.cr(w)
		r.cr(w)
		r.outs(w, "<back>")
		r.cr(w)
	}
	r.documentMatter = node.Matter
}

func (r *Renderer) headingEnter(w io.Writer, heading *ast.Heading) {
	var attrs []string
	tag := "<section"
	if heading.Special != nil {
		tag = "<note"
		if isAbstract(heading.Special) {
			tag = "<abstract"
		}
	}

	r.cr(w)
	r.outTag(w, tag, attrs)
}

func (r *Renderer) headingExit(w io.Writer, heading *ast.Heading) {
	r.cr(w)
}

func (r *Renderer) heading(w io.Writer, node *ast.Heading, entering bool) {
	if !entering {
		r.headingExit(w, node)
		return
	}

	r.sectionClose(w)
	r.section = node
	r.headingEnter(w, node)
}

func (r *Renderer) citation(w io.Writer, node *ast.Citation, entering bool) {
	if !entering {
		return
	}
	for i, c := range node.Destination {
		if node.Type[i] == ast.CitationTypeSuppressed {
			continue
		}

		attr := []string{fmt.Sprintf(`target="%s"`, c)}
		r.outTag(w, "<xref", attr)
		r.outs(w, "</xref>")
	}
}

func (r *Renderer) paragraphEnter(w io.Writer, para *ast.Paragraph) {
	tag := tagWithAttributes("<t", html.BlockAttrs(para))
	r.outs(w, tag)
}

func (r *Renderer) paragraphExit(w io.Writer, para *ast.Paragraph) {
	r.outs(w, "</t>")
	r.cr(w)
}

func (r *Renderer) paragraph(w io.Writer, para *ast.Paragraph, entering bool) {
	if entering {
		r.paragraphEnter(w, para)
	} else {
		r.paragraphExit(w, para)
	}
}

func (r *Renderer) listEnter(w io.Writer, nodeData *ast.List) {
	// TODO: attrs don't seem to be set. Check what is the problem and fix upstream as well.
	var attrs []string

	if nodeData.IsFootnotesList {
		r.outs(w, "\n<div class=\"footnotes\">\n\n")
		r.cr(w)
	}
	r.cr(w)

	openTag := "<ul"
	if nodeData.ListFlags&ast.ListTypeOrdered != 0 {
		if nodeData.Start > 0 {
			attrs = append(attrs, fmt.Sprintf(`start="%d"`, nodeData.Start))
		}
		openTag = "<ol"
	}
	if nodeData.ListFlags&ast.ListTypeDefinition != 0 {
		openTag = "<dl"
	}
	attrs = append(attrs, html.BlockAttrs(nodeData)...)
	r.outTag(w, openTag, attrs)
	r.cr(w)
}

func (r *Renderer) listExit(w io.Writer, list *ast.List) {
	closeTag := "</ul>"
	if list.ListFlags&ast.ListTypeOrdered != 0 {
		closeTag = "</ol>"
	}
	if list.ListFlags&ast.ListTypeDefinition != 0 {
		closeTag = "</dl>"
	}
	r.outs(w, closeTag)

	//cr(w)
	//if node.parent.Type != Item {
	//	cr(w)
	//}
	parent := list.Parent
	switch parent.(type) {
	case *ast.ListItem:
		if ast.GetNextNode(list) != nil {
			r.cr(w)
		}
	case *ast.Document, *ast.BlockQuote, *ast.Aside:
		r.cr(w)
	}

	if list.IsFootnotesList {
		r.outs(w, "\n</div>\n")
	}
}

func (r *Renderer) list(w io.Writer, list *ast.List, entering bool) {
	if entering {
		r.listEnter(w, list)
	} else {
		r.listExit(w, list)
	}
}

func (r *Renderer) listItemEnter(w io.Writer, listItem *ast.ListItem) {
	if listItem.RefLink != nil {
		/*
			TODO
				slug := slugify(listItem.RefLink)
				r.outs(w, footnoteItem(r.opts.FootnoteAnchorPrefix, slug))
		*/
		return
	}

	openTag := "<li>"
	if listItem.ListFlags&ast.ListTypeDefinition != 0 {
		openTag = "<dd>"
	}
	if listItem.ListFlags&ast.ListTypeTerm != 0 {
		openTag = "<dt>"
	}
	r.outs(w, openTag)
}

func (r *Renderer) listItemExit(w io.Writer, listItem *ast.ListItem) {
	closeTag := "</li>"
	if listItem.ListFlags&ast.ListTypeDefinition != 0 {
		closeTag = "</dd>"
	}
	if listItem.ListFlags&ast.ListTypeTerm != 0 {
		closeTag = "</dt>"
	}
	r.outs(w, closeTag)
	r.cr(w)
}

func (r *Renderer) listItem(w io.Writer, listItem *ast.ListItem, entering bool) {
	if entering {
		r.listItemEnter(w, listItem)
	} else {
		r.listItemExit(w, listItem)
	}
}

func (r *Renderer) codeBlock(w io.Writer, codeBlock *ast.CodeBlock) {
	var attrs []string
	attrs = appendLanguageAttr(attrs, codeBlock.Info)
	attrs = append(attrs, html.BlockAttrs(codeBlock)...)

	name := "artwork"
	if codeBlock.Info != nil {
		name = "sourcecode"
	}

	r.cr(w)
	r.outTag(w, "<"+name, attrs)
	if r.opts.Comments != nil {
		r.EscapeHTMLCallouts(w, codeBlock.Literal)
	} else {
		html.EscapeHTML(w, codeBlock.Literal)
	}
	r.outs(w, "</"+name+">")
	r.cr(w)
}

func (r *Renderer) tableCell(w io.Writer, tableCell *ast.TableCell, entering bool) {
	if !entering {
		r.outOneOf(w, tableCell.IsHeader, "</th>", "</td>")
		r.cr(w)
		return
	}

	// entering
	var attrs []string
	openTag := "<td"
	if tableCell.IsHeader {
		openTag = "<th"
	}
	align := tableCell.Align.String()
	if align != "" {
		attrs = append(attrs, fmt.Sprintf(`align="%s"`, align))
	}
	if ast.GetPrevNode(tableCell) == nil {
		r.cr(w)
	}
	r.outTag(w, openTag, attrs)
}

func (r *Renderer) tableBody(w io.Writer, node *ast.TableBody, entering bool) {
	r.outOneOfCr(w, entering, "<tbody>", "</tbody>")
}

func (r *Renderer) htmlSpan(w io.Writer, span *ast.HTMLSpan) {
	if r.opts.Flags&SkipHTML == 0 {
		r.out(w, span.Literal)
	}
}

func (r *Renderer) callout(w io.Writer, callout *ast.Callout) {
	r.outs(w, "<em>")
	r.out(w, callout.ID)
	r.outs(w, "</em>")
}

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
		r.titleBlock(w, node)
	case *mast.References:
		r.references(w, node, entering)
	case *mast.Reference:
		r.reference(w, node)
	case *ast.Text:
		r.text(w, node)
	case *ast.Softbreak:
		r.cr(w)
	case *ast.Hardbreak:
		r.hardBreak(w, node)
	case *ast.Callout:
		r.callout(w, node)
	case *ast.Emph:
		r.outOneOf(w, entering, "<em>", "</em>")
		r.outOneOf(w, entering, "<em>", "</em>")
	case *ast.Strong:
		r.strong(w, node, entering)
	case *ast.Del:
		r.outOneOf(w, entering, "<del>", "</del>")
	case *ast.Citation:
		r.citation(w, node, entering)
	case *ast.DocumentMatter:
		if entering {
			r.matter(w, node)
		}
	case *ast.Heading:
		r.heading(w, node, entering)
	case *ast.Paragraph:
		r.paragraph(w, node, entering)
	case *ast.HTMLSpan:
		r.htmlSpan(w, node) // only html comments are allowed.
	case *ast.HTMLBlock:
		// discard; we use these only for <references>.
	case *ast.List:
		r.list(w, node, entering)
	case *ast.ListItem:
		r.listItem(w, node, entering)
	case *ast.CodeBlock:
		r.codeBlock(w, node)
	case *ast.Caption:
		r.outOneOf(w, entering, "<name>", "</name>")
	case *ast.CaptionFigure:
		r.outOneOf(w, entering, "<figure>", "</figure>")
	case *ast.Table:
		tag := tagWithAttributes("<table", html.BlockAttrs(node))
		r.outOneOfCr(w, entering, tag, "</table>")
	case *ast.TableCell:
		r.tableCell(w, node, entering)
	case *ast.TableHead:
		r.outOneOfCr(w, entering, "<thead>", "</thead>")
	case *ast.TableBody:
		r.tableBody(w, node, entering)
	case *ast.TableRow:
		r.outOneOfCr(w, entering, "<tr>", "</tr>")
	case *ast.BlockQuote:
		tag := tagWithAttributes("<blockquote", html.BlockAttrs(node))
		r.outOneOfCr(w, entering, tag, "</blockquote>")
	case *ast.Aside:
		tag := tagWithAttributes("<aside", html.BlockAttrs(node))
		r.outOneOfCr(w, entering, tag, "</aside>")
	default:
		panic(fmt.Sprintf("Unknown node %T", node))
	}
	return ast.GoToNext
}

// RenderHeader writes HTML document preamble and TOC if requested.
func (r *Renderer) RenderHeader(w io.Writer, ast ast.Node) {
	if r.opts.Flags&XMLFragment != 0 {
		return
	}

	r.writeDocumentHeader(w)
}

// RenderFooter writes HTML document footer.
func (r *Renderer) RenderFooter(w io.Writer, _ ast.Node) {
	r.sectionClose(w)
	r.section = nil

	switch r.documentMatter {
	case ast.DocumentMatterFront:
		r.outs(w, "\n</front>\n")
	case ast.DocumentMatterMain:
		r.outs(w, "\n</middle>\n")
	case ast.DocumentMatterBack:
		r.outs(w, "\n</back>\n")
	}

	if r.opts.Flags&XMLFragment != 0 {
		return
	}

	io.WriteString(w, "\n</rfc>\n")
}

func (r *Renderer) writeDocumentHeader(w io.Writer) {
	r.outs(w, `<?xml version="1.0" encoding="utf-8"?>`)
	r.cr(w)
	r.outs(w, `<meta name="GENERATOR" content="github.com/mmarkdown/mmark markdown processor for Go">`)
	r.cr(w)
	r.outs(w, `<meta charset="utf-8">`)
	r.cr(w)
}

func tagWithAttributes(name string, attrs []string) string {
	s := name
	if len(attrs) > 0 {
		s += " " + strings.Join(attrs, " ")
	}
	return s + ">"
}
