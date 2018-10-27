package xml2

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/render/xml"
)

// Flags control optional behavior of XML2 renderer.
type Flags int

// HTML renderer configuration options.
const (
	FlagsNone   Flags = 0
	XMLFragment Flags = 1 << iota // Don't generate a complete XML document
	SkipHTML                      // Skip preformatted HTML blocks - skips comments
	SkipImages                    // Skip embedded images

	CommonFlags Flags = FlagsNone
)

// RendererOptions is a collection of supplementary parameters tweaking
// the behavior of various parts of XML2 renderer.
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

	// Generator is a comment that is inserted in the generated XML to show what rendered it.
	Generator string
}

// Renderer implements Renderer interface for IETF XMLv2 output. See RFC 7941.
type Renderer struct {
	opts RendererOptions

	documentMatter ast.DocumentMatters // keep track of front/main/back matter
	section        *ast.Heading        // current open section
	title          bool                // did we output a title block
	filter         mast.FilterFunc     // filter for attributes.

	// Track heading IDs to prevent ID collision in a single generation.
	headingIDs map[string]int
}

var filterFunc mast.FilterFunc = func(s string) bool {
	switch s {
	case "id": // will translate to anchor so OK.
		return true
	case "class": // there are no classes
		return false
	case "empty": // newer attributes from RFC 7991
		return false
	}

	// l33t data- HTML5 attributes
	if strings.HasPrefix(s, "data-") {
		return false
	}

	return true
}

// NewRenderer creates and configures an Renderer object, which satisfies the Renderer interface.
func NewRenderer(opts RendererOptions) *Renderer {
	html.IDTag = "anchor"
	if opts.Generator == "" {
		opts.Generator = xml.Generator
	}
	return &Renderer{opts: opts, headingIDs: make(map[string]int), filter: filterFunc}
}

func (r *Renderer) text(w io.Writer, text *ast.Text) {
	if _, parentIsLink := text.Parent.(*ast.Link); parentIsLink {
		r.out(w, text.Literal)
		return
	}

	if heading, parentIsHeading := text.Parent.(*ast.Heading); parentIsHeading {
		if heading.IsSpecial && xml.IsAbstract(heading.Literal) {
			return
		}
	}

	html.EscapeHTML(w, text.Literal)
}

func (r *Renderer) hardBreak(w io.Writer, node *ast.Hardbreak) {
	r.outs(w, "<vspace />")
	r.cr(w)
}

func (r *Renderer) strong(w io.Writer, node *ast.Strong, entering bool) {
	// *iff* we have a text node as a child *and* that text is 2119, we output bcp14 tags, otherwise just string.
	text := ast.GetFirstChild(node)
	if t, ok := text.(*ast.Text); ok {
		if xml.Is2119(t.Literal) {
			// out as-is.
			r.outOneOf(w, entering, "", "")
			return
		}
	}

	if _, isCaption := node.GetParent().(*ast.Caption); isCaption {
		r.outOneOf(w, entering, "", "")
		return
	}

	r.outOneOf(w, entering, `<spanx style="strong">`, "</spanx>")
}

func (r *Renderer) matter(w io.Writer, node *ast.DocumentMatter) {
	r.sectionClose(w, nil)

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
	r.cr(w)

	tag := "<section"
	if heading.IsSpecial {
		tag = "<note"
		if xml.IsAbstract(heading.Literal) {
			tag = "<abstract>"
			r.outs(w, tag)
			return
		}
	}

	mast.AttributeInit(heading)
	if mast.Attribute(heading, "id") == nil && heading.HeadingID != "" {
		id := r.ensureUniqueHeadingID(heading.HeadingID)
		mast.SetAttribute(heading, "id", []byte(id))
	}

	// If we want to support block level attributes here, it will clash with the
	// title= attribute that is outed in text() - and thus later.
	r.outs(w, tag)
	r.outAttr(w, html.BlockAttrs(heading))
	r.outs(w, ` title="`)
}

func (r *Renderer) headingExit(w io.Writer, heading *ast.Heading) {
	if heading.IsSpecial && xml.IsAbstract(heading.Literal) {
		return
	}

	r.outs(w, `">`)
	r.cr(w)
}

func (r *Renderer) heading(w io.Writer, node *ast.Heading, entering bool) {
	if !entering {
		r.headingExit(w, node)
		return
	}

	r.sectionClose(w, node)

	r.headingEnter(w, node)
}

var rule = strings.Repeat("-", 60)

func (r *Renderer) horizontalRule(w io.Writer, node *ast.HorizontalRule) {
	if _, ok := node.Parent.(*ast.ListItem); ok {
		r.outs(w, "<vspace/>")
		r.outs(w, rule)
		r.outs(w, "<vspace/>")
	} else {
		r.outs(w, "<t>")
		r.outs(w, rule)
		r.outs(w, "</t>")
	}
	r.cr(w)
}

func (r *Renderer) citation(w io.Writer, node *ast.Citation, entering bool) {
	if !entering {
		return
	}

	for i, c := range node.Destination {
		if node.Type[i] == ast.CitationTypeSuppressed {
			continue
		}

		r.outTag(w, "<xref", []string{fmt.Sprintf(`target="%s"`, c)})
		r.outs(w, "</xref>")
	}
}

func (r *Renderer) paragraphEnter(w io.Writer, para *ast.Paragraph) {
	// Skip outputting </t> in lists and in caption figures.
	if p, ok := para.Parent.(*ast.ListItem); ok {
		// Fake multiple paragraphs by inserting a hard break.
		if len(p.GetChildren()) > 1 {
			first := ast.GetFirstChild(para.Parent)
			if first != para {
				r.hardBreak(w, &ast.Hardbreak{})
			}
		}
		return
	}
	if _, ok := para.Parent.(*ast.CaptionFigure); ok {
		return
	}

	tag := tagWithAttributes("<t", html.BlockAttrs(para))
	r.outs(w, tag)
}

func (r *Renderer) paragraphExit(w io.Writer, para *ast.Paragraph) {
	// Skip outputting </t> in lists and in caption figures.
	if _, ok := para.Parent.(*ast.ListItem); ok {
		return
	}
	if _, ok := para.Parent.(*ast.CaptionFigure); ok {
		return
	}

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
	if nodeData.IsFootnotesList {
		return
	}
	r.cr(w)

	openTag := "<list"

	// if there is a block level attribute with style, we shouldn't use ours.
	if mast.Attribute(nodeData, "style") == nil {

		mast.AttributeInit(nodeData)

		mast.SetAttribute(nodeData, "style", []byte("symbols"))
		if nodeData.ListFlags&ast.ListTypeOrdered != 0 {
			mast.SetAttribute(nodeData, "style", []byte("numbers"))
			if nodeData.Start > 0 {
				log.Printf("Attribute \"start\" not supported for list style=\"numbers\"")
			}
		}
		if nodeData.ListFlags&ast.ListTypeDefinition != 0 {
			mast.SetAttribute(nodeData, "style", []byte("hanging"))
		}
	}

	r.outTag(w, openTag, html.BlockAttrs(nodeData))
	r.cr(w)
}

func (r *Renderer) listExit(w io.Writer, list *ast.List) {
	if list.IsFootnotesList {
		return
	}
	closeTag := "</list>"
	if list.ListFlags&ast.ListTypeOrdered != 0 {
		//closeTag = "</ol>"
	}
	if list.ListFlags&ast.ListTypeDefinition != 0 {
		//closeTag = "</dl>"
	}
	r.outs(w, closeTag)

	parent := list.Parent
	switch parent.(type) {
	case *ast.ListItem:
		if ast.GetNextNode(list) != nil {
			r.cr(w)
		}
	case *ast.Document, *ast.BlockQuote, *ast.Aside:
		r.cr(w)
	}
}

func (r *Renderer) list(w io.Writer, list *ast.List, entering bool) {
	// need to be wrapped in a paragraph, except when we're already in a list.
	_, parentIsList := list.Parent.(*ast.ListItem)
	if entering {
		if !parentIsList {
			r.paragraphEnter(w, &ast.Paragraph{})
		}
		r.listEnter(w, list)
	} else {
		r.listExit(w, list)
		if !parentIsList {
			r.paragraphExit(w, &ast.Paragraph{})
		}
	}
}

func (r *Renderer) listItemEnter(w io.Writer, listItem *ast.ListItem) {
	if listItem.RefLink != nil { // footnotes
		return
	}

	openTag := "<t>"
	if listItem.ListFlags&ast.ListTypeDefinition != 0 {
		openTag = "<vspace />"
	}
	if listItem.ListFlags&ast.ListTypeTerm != 0 {
		openTag = "<t hangText=\""
	}
	r.outs(w, openTag)
}

func (r *Renderer) listItemExit(w io.Writer, listItem *ast.ListItem) {
	if listItem.RefLink != nil {
		return
	}

	closeTag := "</t>"
	if listItem.ListFlags&ast.ListTypeTerm != 0 {
		closeTag = `">`
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
	mast.AttributeInit(codeBlock)
	appendLanguageAttr(codeBlock, codeBlock.Info)

	r.cr(w)
	_, inFigure := codeBlock.Parent.(*ast.CaptionFigure)
	if inFigure {
		// Drop anchor for now, but need to figure out what to allow here.
		mast.DeleteAttribute(codeBlock, "id")
		r.outTag(w, "<artwork", html.BlockAttrs(codeBlock))
	} else {
		typ := mast.Attribute(codeBlock, "type") // only valid on artwork
		mast.DeleteAttribute(codeBlock, "type")
		r.outTag(w, "<figure", html.BlockAttrs(codeBlock))
		mast.DeleteAttribute(codeBlock, "id")
		if typ != nil {
			mast.SetAttribute(codeBlock, "type", typ)
		}
		r.outTag(w, "<artwork", html.BlockAttrs(codeBlock))
	}

	if r.opts.Comments != nil {
		xml.EscapeHTMLCallouts(w, codeBlock.Literal, r.opts.Comments)
	} else {
		html.EscapeHTML(w, codeBlock.Literal)
	}
	if inFigure {
		r.outs(w, "</artwork>")
	} else {
		r.outs(w, "</artwork></figure>\n")
	}
	r.cr(w)
}

func (r *Renderer) tableCell(w io.Writer, tableCell *ast.TableCell, entering bool) {
	if !entering {
		r.outOneOf(w, tableCell.IsHeader, "</ttcol>", "</c>")
		r.cr(w)
		return
	}

	// entering
	mast.AttributeInit(tableCell)
	openTag := "<c"
	if tableCell.IsHeader {
		openTag = "<ttcol"
		align := tableCell.Align.String()
		if align != "" {
			mast.SetAttribute(tableCell, "align", []byte(align))
		}
	}
	if ast.GetPrevNode(tableCell) == nil {
		r.cr(w)
	}
	r.outTag(w, openTag, html.BlockAttrs(tableCell))
}

func (r *Renderer) tableBody(w io.Writer, node *ast.TableBody, entering bool) {
	r.outOneOfCr(w, entering, "", "")
}

func (r *Renderer) htmlSpan(w io.Writer, span *ast.HTMLSpan) {
	if r.opts.Flags&SkipHTML == 0 {
		html.EscapeHTML(w, span.Literal)
	}
}

func (r *Renderer) callout(w io.Writer, callout *ast.Callout) {
	r.outs(w, `<spanx style="emph">`)
	r.out(w, callout.ID)
	r.outs(w, "</spanx>")
}

func (r *Renderer) crossReference(w io.Writer, cr *ast.CrossReference, entering bool) {
	if isHangText(cr) {
		if entering {
			w.Write(cr.Destination)
		}
		return
	}

	if entering {
		r.outTag(w, "<xref", []string{"target=\"" + string(cr.Destination) + "\""})
		return
	}
	r.outs(w, "</xref>")
}

func (r *Renderer) index(w io.Writer, index *ast.Index) {
	r.outs(w, "<iref")
	r.outs(w, " item=\"")
	html.EscapeHTML(w, index.Item)
	r.outs(w, "\"")
	if index.Primary {
		r.outs(w, ` primary="true"`)
	}
	if len(index.Subitem) != 0 {
		r.outs(w, " subitem=\"")
		html.EscapeHTML(w, index.Subitem)
		r.outs(w, "\"")
	}
	r.outs(w, "/>")
}

func (r *Renderer) link(w io.Writer, link *ast.Link, entering bool) {
	if link.Footnote != nil {
		return
	}

	if !entering {
		r.outs(w, `</eref>`)
		return
	}

	if isHangText(link) {
		w.Write(link.Destination)
		return
	}

	r.outs(w, "<eref")
	r.outs(w, " target=\"")
	html.EscapeHTML(w, link.Destination)
	r.outs(w, `">`)
}

func (r *Renderer) image(w io.Writer, node *ast.Image, entering bool) {
	if entering {
		r.imageEnter(w, node)
	} else {
		r.imageExit(w, node)
	}
}

func (r *Renderer) imageEnter(w io.Writer, image *ast.Image) {
	r.outs(w, "<artwork>\n")
	html.EscapeHTML(w, image.Destination)
	r.outs(w, ` `)
}

func (r *Renderer) imageExit(w io.Writer, image *ast.Image) {
	if image.Title != nil {
		r.outs(w, ` "`)
		html.EscapeHTML(w, image.Title)
		r.outs(w, `"`)
	}
	r.outs(w, "</artwork>\n")
}

func (r *Renderer) code(w io.Writer, node *ast.Code) {
	if isHangText(node) {
		html.EscapeHTML(w, node.Literal)
		return
	}
	if _, isCaption := node.GetParent().(*ast.Caption); isCaption {
		html.EscapeHTML(w, node.Literal)
		return
	}

	r.outs(w, `<spanx style="verb">`)
	html.EscapeHTML(w, node.Literal)
	r.outs(w, "</spanx>")
}

func (r *Renderer) mathBlock(w io.Writer, mathBlock *ast.MathBlock) {
	r.outs(w, `<figure><artwork type="math">`+"\n")
	if r.opts.Comments != nil {
		xml.EscapeHTMLCallouts(w, mathBlock.Literal, r.opts.Comments)
	} else {
		html.EscapeHTML(w, mathBlock.Literal)
	}
	r.outs(w, `</artwork></figure>`+"\n")
	r.cr(w)
}

func (r *Renderer) captionFigure(w io.Writer, captionFigure *ast.CaptionFigure, entering bool) {
	for _, child := range captionFigure.GetChildren() {
		if _, ok := child.(*ast.Table); ok {
			return
		}
		if _, ok := child.(*ast.BlockQuote); ok {
			return
		}
	}

	if !entering {
		r.outs(w, "</figure>\n")
		return
	}
	if captionFigure.HeadingID != "" {
		mast.AttributeInit(captionFigure)
		captionFigure.Attribute.ID = []byte(captionFigure.HeadingID)
	}

	r.outs(w, "<figure")
	r.outAttr(w, html.BlockAttrs(captionFigure))

	// Now render the caption and then *remove* it from the tree.
	for _, child := range captionFigure.GetChildren() {
		if caption, ok := child.(*ast.Caption); ok {
			r.outs(w, ` title="`)
			ast.WalkFunc(caption, func(node ast.Node, entering bool) ast.WalkStatus {
				return r.RenderNode(w, node, entering)
			})
			r.outs(w, `">`)

			ast.RemoveFromTree(caption)
			return
		}
	}
	// Still here? Close tag.
	r.outs(w, `>`)
}

func (r *Renderer) table(w io.Writer, tab *ast.Table, entering bool) {
	if !entering {
		r.outs(w, "</texttable>")
		return
	}
	captionFigure, inFigure := tab.Parent.(*ast.CaptionFigure)
	if inFigure {
		if captionFigure.HeadingID != "" {
			mast.AttributeInit(tab)
			tab.Attribute.ID = []byte(captionFigure.HeadingID)
		}
	}

	r.outs(w, "<texttable")
	r.outAttr(w, html.BlockAttrs(tab))
	// Now render the caption if our parent is a ast.CaptionFigure
	// and then *remove* it from the tree.
	if !inFigure {
		r.outs(w, `>`)
		return
	}

	r.outs(w, ` title="`)

	for _, child := range captionFigure.GetChildren() {
		if caption, ok := child.(*ast.Caption); ok {
			ast.WalkFunc(caption, func(node ast.Node, entering bool) ast.WalkStatus {
				return r.RenderNode(w, node, entering)
			})
			r.outs(w, `">`)

			ast.RemoveFromTree(caption)
			return
		}
	}
	// Still here? Close tag.
	r.outs(w, `>`)
}

func (r *Renderer) blockQuote(w io.Writer, block *ast.BlockQuote, entering bool) {
	if !entering {
		return
	}

	// Fake a list. TODO(miek): list in list checks, see? Make fake parent??
	list := &ast.List{}
	list.Attribute = block.Attribute
	if list.Attribute == nil {
		list.Attribute = &ast.Attribute{Attrs: make(map[string][]byte)}
	}
	list.Attribute.Attrs["style"] = []byte("empty")

	listItem := &ast.ListItem{}
	mast.MoveChildren(listItem, block)
	ast.AppendChild(list, listItem)

	if captionFigure, ok := block.Parent.(*ast.CaptionFigure); ok {
		for _, child := range captionFigure.GetChildren() {
			if caption, ok := child.(*ast.Caption); ok {
				listItem := &ast.ListItem{}
				mast.MoveChildren(listItem, caption)
				ast.AppendChild(list, listItem)

				ast.RemoveFromTree(caption)
			}
		}
	}

	ast.WalkFunc(list, func(node ast.Node, entering bool) ast.WalkStatus {
		return r.RenderNode(w, node, entering)
	})
}

// RenderNode renders a markdown node to XML.
func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {

	mast.AttributeFilter(node, r.filter)

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
		r.title = true
	case *mast.Bibliography:
		r.bibliography(w, node, entering)
	case *mast.BibliographyItem:
		r.bibliographyItem(w, node)
	case *mast.DocumentIndex, *mast.IndexLetter, *mast.IndexItem, *mast.IndexSubItem, *mast.IndexLink:
		// generated by xml2rfc, do nothing.
	case *ast.Text:
		r.text(w, node)
	case *ast.Softbreak:
		r.cr(w)
	case *ast.Hardbreak:
		r.hardBreak(w, node)
	case *ast.Callout:
		r.callout(w, node)
	case *ast.Emph:
		if isHangText(node) {
			if entering {
				html.EscapeHTML(w, node.Literal)
			}
		} else {
			if _, isCaption := node.GetParent().(*ast.Caption); isCaption {
				r.outOneOf(w, entering, "", "")
			} else {
				r.outOneOf(w, entering, `<spanx style="emph">`, "</spanx>")
			}
		}
	case *ast.Strong:
		r.strong(w, node, entering)
	case *ast.Del:
		// ala strikethrough, just keep the tildes
		r.outOneOf(w, entering, "~", "~")
	case *ast.Citation:
		r.citation(w, node, entering)
	case *ast.DocumentMatter:
		if entering {
			r.matter(w, node)
		}
	case *ast.Heading:
		r.heading(w, node, entering)
	case *ast.HorizontalRule:
		r.horizontalRule(w, node)
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
		// no tags because we are used in attributes, i.e. title=
		// See comment in xml/renderer.go. The same is true here, *but*, because we don't
		// output any tags, the problem does not show up. As a matter of consistency we apply
		// the same (dumb) precaution.
		if len(node.GetChildren()) > 0 {
			r.outOneOf(w, entering, "", "")
		}
	case *ast.CaptionFigure:
		r.captionFigure(w, node, entering)
	case *ast.Table:
		r.table(w, node, entering)
	case *ast.TableCell:
		r.tableCell(w, node, entering)
	case *ast.TableHeader:
		r.outOneOf(w, entering, "", "")
	case *ast.TableBody:
		r.tableBody(w, node, entering)
	case *ast.TableRow:
		r.outOneOf(w, entering, "", "")
	case *ast.TableFooter:
		r.outOneOf(w, entering, "", "")
	case *ast.BlockQuote:
		r.blockQuote(w, node, entering)
	case *ast.Aside:
		// ignore and text render the child text as-is.
	case *ast.CrossReference:
		r.crossReference(w, node, entering)
	case *ast.Index:
		if entering {
			r.index(w, node)
		}
	case *ast.Link:
		r.link(w, node, entering)
	case *ast.Math:
		r.outOneOf(w, true, `<spanx style="verb">`, "</spanx>")
		html.EscapeHTML(w, node.Literal)
		r.outOneOf(w, false, `<spanx style="verb">`, "</spanx>")
	case *ast.Image:
		if r.opts.Flags&SkipImages != 0 {
			return ast.SkipChildren
		}
		r.image(w, node, entering)
	case *ast.Code:
		r.code(w, node)
	case *ast.MathBlock:
		r.mathBlock(w, node)
	case *ast.Subscript:
		r.outOneOf(w, true, "_(", ")")
		if entering {
			html.Escape(w, node.Literal)
		}
		r.outOneOf(w, false, "_(", ")")
	case *ast.Superscript:
		r.outOneOf(w, true, "^(", ")")
		if entering {
			html.Escape(w, node.Literal)
		}
		r.outOneOf(w, false, "^(", ")")
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
	r.sectionClose(w, nil)

	switch r.documentMatter {
	case ast.DocumentMatterFront:
		r.outs(w, "\n</front>\n")
	case ast.DocumentMatterMain:
		r.outs(w, "\n</middle>\n")
	case ast.DocumentMatterBack:
		r.outs(w, "\n</back>\n")
	}

	if r.title {
		io.WriteString(w, "\n</rfc>")
	}
}

func (r *Renderer) writeDocumentHeader(w io.Writer) {
	if r.opts.Flags&XMLFragment != 0 {
		return
	}
	r.outs(w, `<?xml version="1.0" encoding="utf-8"?>`)
	r.cr(w)
	r.outs(w, r.opts.Generator)
	r.cr(w)
	r.outs(w, `<!DOCTYPE rfc SYSTEM 'rfc2629.dtd' []>`)
	r.cr(w)
}

func tagWithAttributes(name string, attrs []string) string {
	s := name
	if len(attrs) > 0 {
		s += " " + strings.Join(attrs, " ")
	}
	return s + ">"
}
