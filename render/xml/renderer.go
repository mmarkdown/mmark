package xml

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/mmarkdown/mmark/mast"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

// Flags control optional behavior of XML3 renderer.
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
// the behavior of various parts of XML renderer.
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

// Renderer implements Renderer interface for IETF XMLv3 output. See RFC 7991.
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
	case "style": // style has been deprecated in 7991
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
		opts.Generator = Generator
	}
	return &Renderer{opts: opts, headingIDs: make(map[string]int), filter: filterFunc}
}

func (r *Renderer) text(w io.Writer, text *ast.Text) {
	if _, parentIsLink := text.Parent.(*ast.Link); parentIsLink {
		//html.EscLink(w, text.Literal)
		r.out(w, text.Literal)
		return
	}

	if heading, parentIsHeading := text.Parent.(*ast.Heading); parentIsHeading {
		if heading.IsSpecial && IsAbstract(heading.Literal) {
			// No <name> when abstract, should not output anything
			// This works because abstract does not contain any markdown, i.e. <em>Abstract</em> would still output the emphesis.
			return
		}
	}
	// each string of unicode chars is outputted between <u> tags.
	uni := 0
	for i, c := range string(text.Literal) {
		if c > unicode.MaxASCII {
			uni += len(string(c))
			continue
		}
		if uni > 0 {
			r.outs(w, `<u format="char-num">`)
			r.outs(w, string(text.Literal)[i-uni:i]) // must we also html escape this??
			r.outs(w, `</u>`)
			uni = 0
		}
		html.EscapeHTML(w, []byte(string(c)))
	}
	// last chars where uni
	if uni > 0 {
		i := len(string(text.Literal))
		r.outs(w, `<u format="char-num">`)
		r.outs(w, string(text.Literal)[i-uni:i])
		r.outs(w, `</u>`)
	}

}

func (r *Renderer) hardBreak(w io.Writer, node *ast.Hardbreak) {
	r.outs(w, "<br />")
	r.cr(w)
}

func (r *Renderer) strong(w io.Writer, node *ast.Strong, entering bool) {
	// *iff* we have a text node as a child *and* that text is 2119, we output bcp14 tags, otherwise just string.
	text := ast.GetFirstChild(node)
	if t, ok := text.(*ast.Text); ok {
		if Is2119(t.Literal) {
			r.outOneOf(w, entering, "<bcp14>", "</bcp14>")
			return
		}
	}

	r.outOneOf(w, entering, "<strong>", "</strong>")
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
	tag := "<section"

	mast.AttributeInit(heading)
	// In XML2 output we can't have an anchor attribute on a note.
	if !heading.IsSpecial {
		if mast.Attribute(heading, "id") == nil && heading.HeadingID != "" {
			id := r.ensureUniqueHeadingID(heading.HeadingID)
			mast.SetAttribute(heading, "id", []byte(id))
		}
	}

	if heading.IsSpecial {
		tag = "<note"
		if IsAbstract(heading.Literal) {
			tag = "<abstract"
		}
	}

	r.cr(w)
	r.outTag(w, tag, html.BlockAttrs(heading))

	if heading.IsSpecial && IsAbstract(heading.Literal) {
		return
	}
	r.outs(w, "<name>")
}

func (r *Renderer) headingExit(w io.Writer, heading *ast.Heading) {
	if heading.IsSpecial && IsAbstract(heading.Literal) {
		r.cr(w)
		return
	}
	r.outs(w, "</name>")
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

func (r *Renderer) citation(w io.Writer, node *ast.Citation, entering bool) {
	if !entering {
		return
	}
	for i, c := range node.Destination {
		if node.Type[i] == ast.CitationTypeSuppressed {
			continue
		}

		attr := []string{fmt.Sprintf(`target="%s"`, c)}
		// Attempt to parse the suffix, supported:
		// 'section N', where N is a number. We're not checking the number, garbage in, garbage out
		if len(node.Suffix) != 0 {
			if bytes.HasPrefix(node.Suffix[0], []byte("section ")) {
				num := node.Suffix[0][len("section "):]
				attr = append(attr, `sectionFormat="bare" relative="#"`)
				attr = append(attr, `section="`+string(num)+`"`)
			}
		}

		r.outTag(w, "<xref", attr)
		r.outs(w, "</xref>")
	}
}

func (r *Renderer) paragraphEnter(w io.Writer, para *ast.Paragraph) {
	if p, ok := para.Parent.(*ast.ListItem); ok {
		if p.ListFlags&ast.ListTypeTerm != 0 {
			return
		}
		if p.ListFlags&ast.ListItemContainsBlock == 0 {
			// No block level elements, don't output a paragraph. However a list in a list
			// doesn't show up as a block level element, so we check for that as well. In that
			// case we do need to output a para tag.
			listInList := false
			ast.WalkFunc(p, func(node ast.Node, entering bool) ast.WalkStatus {
				if _, ok := node.(*ast.List); ok {
					listInList = true
					return ast.Terminate
				}
				return ast.GoToNext
			})
			if !listInList {
				return
			}
		}
	}
	if _, ok := para.Parent.(*ast.CaptionFigure); ok {
		return
	}
	tag := tagWithAttributes("<t", html.BlockAttrs(para))
	r.outs(w, tag)
}

func (r *Renderer) paragraphExit(w io.Writer, para *ast.Paragraph) {
	if p, ok := para.Parent.(*ast.ListItem); ok {
		if p.ListFlags&ast.ListTypeTerm != 0 {
			return
		}
		if p.ListFlags&ast.ListItemContainsBlock == 0 {
			// see comments in paragraphEnter
			listInList := false
			ast.WalkFunc(p, func(node ast.Node, entering bool) ast.WalkStatus {
				if _, ok := node.(*ast.List); ok {
					listInList = true
					return ast.Terminate
				}
				return ast.GoToNext
			})
			if !listInList {
				return
			}
		}
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

	openTag := "<ul"
	mast.AttributeInit(nodeData)
	if nodeData.ListFlags&ast.ListTypeOrdered != 0 {
		if nodeData.Start > 0 {
			mast.SetAttribute(nodeData, "start", []byte(strconv.Itoa(nodeData.Start)))
		}
		openTag = "<ol"
	}
	if nodeData.ListFlags&ast.ListTypeDefinition != 0 {
		openTag = "<dl"
	}
	r.outTag(w, openTag, html.BlockAttrs(nodeData))
	r.cr(w)
}

func (r *Renderer) listExit(w io.Writer, list *ast.List) {
	if list.IsFootnotesList {
		return
	}
	closeTag := "</ul>"
	if list.ListFlags&ast.ListTypeOrdered != 0 {
		closeTag = "</ol>"
	}
	if list.ListFlags&ast.ListTypeDefinition != 0 {
		closeTag = "</dl>"
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
	if listItem.RefLink != nil { //footnotes
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
	if listItem.RefLink != nil {
		return
	}

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
	mast.AttributeInit(codeBlock)
	appendLanguageAttr(codeBlock, codeBlock.Info)

	name := "artwork"
	if codeBlock.Info != nil {
		name = "sourcecode"
	}

	r.cr(w)
	r.outTag(w, "<"+name, html.BlockAttrs(codeBlock))
	if r.opts.Comments != nil {
		EscapeHTMLCallouts(w, codeBlock.Literal, r.opts.Comments)
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
	mast.AttributeInit(tableCell)
	openTag := "<td"
	if tableCell.IsHeader {
		openTag = "<th"
	}
	align := tableCell.Align.String()
	if align != "" {
		mast.SetAttribute(tableCell, "align", []byte(align))
	}
	if colspan := tableCell.ColSpan; colspan > 0 {
		mast.SetAttribute(tableCell, "colspan", []byte(fmt.Sprintf("%d", colspan)))

	}
	if ast.GetPrevNode(tableCell) == nil {
		r.cr(w)
	}
	r.outTag(w, openTag, html.BlockAttrs(tableCell))
}

func (r *Renderer) tableBody(w io.Writer, node *ast.TableBody, entering bool) {
	r.outOneOfCr(w, entering, "<tbody>", "</tbody>")
}

func (r *Renderer) htmlSpan(w io.Writer, span *ast.HTMLSpan) {
	if _, ok := IsComment(span.Literal); ok {
		return
	}
	if IsBr(span.Literal) {
		r.hardBreak(w, &ast.Hardbreak{})
		return
	}

	if r.opts.Flags&SkipHTML == 0 {
		html.EscapeHTML(w, span.Literal)
	}
}

func (r *Renderer) htmlBlock(w io.Writer, block *ast.HTMLBlock) {
	if _, ok := IsComment(block.Literal); ok {
		return
	}
	if r.opts.Flags&SkipHTML == 0 {
		html.EscapeHTML(w, block.Literal)
	}
}

func (r *Renderer) callout(w io.Writer, callout *ast.Callout) {
	r.outs(w, "<em>")
	r.out(w, callout.ID)
	r.outs(w, "</em>")
}

func (r *Renderer) crossReference(w io.Writer, cr *ast.CrossReference, entering bool) {
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
	r.outs(w, "<eref")
	r.outs(w, ` target="`)
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
	dest := image.Destination
	r.outs(w, `<artwork src="`)
	html.EscapeHTML(w, dest)
	r.outs(w, `"`)
	ext := path.Ext(string(dest))
	if len(ext) > 2 {
		// warn if not svn or ascii-art.
		switch ext {
		case ".svg", ".ascii-art":
		default:
			log.Printf("Image extension of %q will likely create errors in XML2RFC", ext)
		}
		r.outs(w, ` type="`)
		r.outs(w, ext[1:])
		r.outs(w, `"`)
	}

	// alt= will be the alt text (which is normal text so rendered by the normal render flow)
	r.outs(w, ` alt="`)
}

func (r *Renderer) imageExit(w io.Writer, image *ast.Image) {
	// where to put image title? Put in the artwork?
	if image.Title != nil {
		r.outs(w, `" name="`)
		html.EscapeHTML(w, image.Title)
	}
	r.outs(w, `"/>`)
}

func (r *Renderer) code(w io.Writer, node *ast.Code) {
	r.outs(w, "<tt>")
	html.EscapeHTML(w, node.Literal)
	r.outs(w, "</tt>")
}

func (r *Renderer) mathBlock(w io.Writer, mathBlock *ast.MathBlock) {
	r.outs(w, `<artwork type="math">`+"\n")
	if r.opts.Comments != nil {
		EscapeHTMLCallouts(w, mathBlock.Literal, r.opts.Comments)
	} else {
		html.EscapeHTML(w, mathBlock.Literal)
	}
	r.outs(w, `</artwork>`)
	r.cr(w)
}

func (r *Renderer) captionFigure(w io.Writer, captionFigure *ast.CaptionFigure, entering bool) {
	// If the captionFigure has a table as child element *don't* output the figure tags, because 7991 is weird.
	// If we have a quoted blockquote it is also wrapped in a figure, which we don't want.
	//
	// To detect an artset, we check the number of images, *and* if the filename referenced in Destination is equal apart from the
	// extensions. We don't care about the extension here, but if we detect this we wrap the lot in a artset.
	base := ""
	artset := false
	for _, child := range captionFigure.GetChildren() {
		if _, ok := child.(*ast.Table); ok {
			return
		}
		if _, ok := child.(*ast.BlockQuote); ok {
			return
		}
		// for figures/images, these are wrapped below a paragraph.
		// TODO: can there be more than 1 paragraph??
		if p, ok := child.(*ast.Paragraph); ok {
			for _, img := range p.GetChildren() {
				x, ok := img.(*ast.Image)
				if !ok {
					continue
				}
				b := strings.TrimSuffix(string(x.Destination), filepath.Ext(string(x.Destination)))
				artset = false
				if b == base {
					artset = true
				}
				base = b
			}
		}
	}

	if !entering {
		if artset {
			r.outs(w, "</artset>\n")
		}
		r.outs(w, "</figure>\n")
		return
	}

	if captionFigure.HeadingID != "" {
		mast.AttributeInit(captionFigure)
		captionFigure.Attribute.ID = []byte(captionFigure.HeadingID)
	}

	r.outs(w, "<figure")
	r.outAttr(w, html.BlockAttrs(captionFigure))
	r.outs(w, ">")

	// Now render the caption and then *remove* it from the tree.
	for _, child := range captionFigure.GetChildren() {
		if caption, ok := child.(*ast.Caption); ok {
			ast.WalkFunc(caption, func(node ast.Node, entering bool) ast.WalkStatus {
				return r.RenderNode(w, node, entering)
			})

			ast.RemoveFromTree(caption)
			break
		}
	}

	// Remove the anchor from the artwork otherwise we have a duplicate ID.
	for _, child := range captionFigure.GetChildren() {
		if artwork, ok := child.(*ast.CodeBlock); ok {
			mast.DeleteAttribute(artwork, "id")
		}
	}
	if artset {
		r.outs(w, "<artset>\n")
	}

}

func (r *Renderer) table(w io.Writer, tab *ast.Table, entering bool) {
	if !entering {
		r.outs(w, "</table>")
		return
	}
	captionFigure, isCaptionFigure := tab.Parent.(*ast.CaptionFigure)
	if isCaptionFigure && captionFigure.HeadingID != "" {
		mast.AttributeInit(tab)
		tab.Attribute.ID = []byte(captionFigure.HeadingID)
	}

	tag := tagWithAttributes("<table", html.BlockAttrs(tab))
	r.outs(w, tag)

	// Now render the caption if our parent is a ast.CaptionFigure
	// and then *remove* it from the tree.
	if !isCaptionFigure {
		return
	}
	for _, child := range captionFigure.GetChildren() {
		if caption, ok := child.(*ast.Caption); ok {
			ast.WalkFunc(caption, func(node ast.Node, entering bool) ast.WalkStatus {
				return r.RenderNode(w, node, entering)
			})

			ast.RemoveFromTree(caption)
			break
		}
	}
}

func (r *Renderer) blockQuote(w io.Writer, block *ast.BlockQuote, entering bool) {
	if r.section != nil && r.section.IsSpecial {
		// XML2RFC doesn't like blocklevel elements in special section, this should be done in other
		// places, but I'm using this in learning go and I want that to translate cleanly to txt with xml2rfc.
		return
	}

	if !entering {
		r.outs(w, "</blockquote>")
		return
	}
	captionFigure, isCaptionFigure := block.Parent.(*ast.CaptionFigure)
	if isCaptionFigure && captionFigure.HeadingID != "" {
		mast.AttributeInit(block)
		block.Attribute.ID = []byte(captionFigure.HeadingID)
	}

	r.outs(w, "<blockquote")
	r.outAttr(w, html.BlockAttrs(block))
	defer r.outs(w, ">")

	// Now render the caption if our parent is a ast.CaptionFigure
	// and then *remove* it from the tree.
	if !isCaptionFigure {
		return
	}
	for _, child := range captionFigure.GetChildren() {
		if caption, ok := child.(*ast.Caption); ok {
			// So we can't render this as-is, because we're putting is in a attribute
			// so we should loose the tags. Hence we render each child separate. This may
			// still create tags, which is wrong, but up to the user.

			if len(caption.GetChildren()) > 0 {
				r.outs(w, ` quotedFrom="`)
			}
			for _, child1 := range caption.GetChildren() {
				ast.WalkFunc(child1, func(node ast.Node, entering bool) ast.WalkStatus {
					return r.RenderNode(w, node, entering)
				})
			}
			r.outs(w, `"`) // closes quotedFrom

			ast.RemoveFromTree(caption)
			break
		}
	}
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
	case *mast.Authors:
		// ignore
	case *mast.Bibliography:
		r.bibliography(w, node, entering)
	case *mast.BibliographyItem:
		r.bibliographyItem(w, node)
	case *mast.DocumentIndex, *mast.IndexLetter, *mast.IndexItem, *mast.IndexSubItem, *mast.IndexLink:
		// generated by xml2rfc, do nothing.
	case *mast.ReferenceBlock:
		// skip, added and done by AddBibliography
	case *ast.Text:
		r.text(w, node)
	case *ast.Softbreak:
		r.cr(w)
	case *ast.Hardbreak:
		r.hardBreak(w, node)
	case *ast.NonBlockingSpace:
		// skip
	case *ast.Callout:
		r.callout(w, node)
	case *ast.Emph:
		r.outOneOf(w, entering, "<em>", "</em>")
	case *ast.Strong:
		r.strong(w, node, entering)
	case *ast.Del:
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
		/* not supported */
	case *ast.Paragraph:
		r.paragraph(w, node, entering)
	case *ast.HTMLSpan:
		r.htmlSpan(w, node)
	case *ast.HTMLBlock:
		r.htmlBlock(w, node)
	case *ast.List:
		r.list(w, node, entering)
	case *ast.ListItem:
		r.listItem(w, node, entering)
	case *ast.CodeBlock:
		r.codeBlock(w, node)
	case *ast.Caption:
		// We do some funky node re-ordering for the caption so it is rendered in the correct
		// spot. For some reason -- even when we call ast.RemoveFromTree -- we still end up
		// here, rendering a caption. I have no idea (yet) why this is. The only good thing is
		// that in these cases there are no children of the caption. So we check that and refuse
		// to render anything.
		if len(node.GetChildren()) > 0 {
			r.outOneOf(w, entering, "<name>", "</name>")
		}
	case *ast.CaptionFigure:
		r.captionFigure(w, node, entering)
	case *ast.Table:
		r.table(w, node, entering)
	case *ast.TableCell:
		r.tableCell(w, node, entering)
	case *ast.TableHeader:
		r.outOneOfCr(w, entering, "<thead>", "</thead>")
	case *ast.TableBody:
		r.tableBody(w, node, entering)
	case *ast.TableRow:
		r.outOneOfCr(w, entering, "<tr>", "</tr>")
	case *ast.TableFooter:
		r.outOneOfCr(w, entering, "<tfoot>", "</tfoot>")
	case *ast.BlockQuote:
		r.blockQuote(w, node, entering)
	case *ast.Aside:
		tag := tagWithAttributes("<aside", html.BlockAttrs(node))
		r.outOneOfCr(w, entering, tag, "</aside>")
	case *ast.CrossReference:
		r.crossReference(w, node, entering)
	case *ast.Index:
		if entering {
			r.index(w, node)
		}
	case *ast.Link:
		r.link(w, node, entering)
	case *ast.Math:
		r.outOneOf(w, true, "<tt>", "</tt>")
		html.EscapeHTML(w, node.Literal)
		r.outOneOf(w, false, "<tt>", "</tt>")
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
		r.outOneOf(w, true, "<sub>", "</sub>")
		if entering {
			html.Escape(w, node.Literal)
		}
		r.outOneOf(w, false, "<sub>", "</sub>")
	case *ast.Superscript:
		r.outOneOf(w, true, "<sup>", "</sup>")
		if entering {
			html.Escape(w, node.Literal)
		}
		r.outOneOf(w, false, "<sup>", "</sup>")
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
}

func tagWithAttributes(name string, attrs []string) string {
	s := name
	if len(attrs) > 0 {
		s += " " + strings.Join(attrs, " ")
	}
	return s + ">"
}

const Generator = `<!-- name="GENERATOR" content="github.com/mmarkdown/mmark Mmark Markdown Processor - mmark.miek.nl" -->`
