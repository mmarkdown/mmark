// The package man outputs man pages from mmmark markdown.
package man

// Lots of code copied from https://github.com/cpuguy83/go-md2man, but adapated to mmark
// and made to support mmark features.

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/mmarkdown/mmark/mast"
)

// Flags control optional behavior of Markdown renderer.
type Flags int

// HTML renderer configuration options.
const (
	FlagsNone Flags = 0

	CommonFlags Flags = FlagsNone
)

// RendererOptions is a collection of supplementary parameters tweaking
// the behavior of various parts of Markdown renderer.
type RendererOptions struct {
	Flags Flags // Flags allow customizing this renderer's behavior

	// if set, called at the start of RenderNode(). Allows replacing rendering of some nodes
	RenderNodeHook html.RenderNodeFunc
}

// Renderer implements Renderer interface for Markdown output.
type Renderer struct {
	opts RendererOptions

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
	// TODO
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
	if i == 0 {
		// maybe error later
		i = len(node.Title)
	}
	if i > 0 {
		r.outs(w, fmt.Sprintf(".TH %q", strings.ToUpper(node.Title[:i-1])))
		r.outs(w, fmt.Sprintf(" %q", node.Title[i:]))
	}
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
		default:
			r.outs(w, "\n.SS ")
		}
	}
}

func (r *Renderer) citation(w io.Writer, node *ast.Citation, entering bool) {
	r.outs(w, "[@")
	for i, dest := range node.Destination {
		if i > 0 {
			r.outs(w, ", ")
		}
		switch node.Type[i] {
		case ast.CitationTypeInformative:
			// skip outputting ? as it's the default
		case ast.CitationTypeNormative:
			r.outs(w, "!")
		case ast.CitationTypeSuppressed:
			r.outs(w, "-")
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
	// needs other types of lists as well, now just the simple one
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
			r.outs(w, fmt.Sprintf(".IP %d\\. 4\n", i+1))

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

func (r *Renderer) table(w io.Writer, tab *ast.Table, entering bool) {}

func (r *Renderer) tableRow(w io.Writer, tableRow *ast.TableRow, entering bool) {}

func (r *Renderer) tableCell(w io.Writer, tableCell *ast.TableCell, entering bool) {}

func (r *Renderer) htmlSpan(w io.Writer, span *ast.HTMLSpan) {}

func (r *Renderer) crossReference(w io.Writer, cr *ast.CrossReference, entering bool) {}

func (r *Renderer) index(w io.Writer, index *ast.Index, entering bool) {}

func (r *Renderer) link(w io.Writer, link *ast.Link, entering bool) {
	// !entering so the URL comes after the link text.
	if !entering {
		r.outs(w, "\n\\[la]")
		r.out(w, link.Destination)
		r.outs(w, "\\[ra]")
	}
}

func (r *Renderer) image(w io.Writer, node *ast.Image, entering bool) {}

func (r *Renderer) mathBlock(w io.Writer, mathBlock *ast.MathBlock, entering bool) {
}

func (r *Renderer) captionFigure(w io.Writer, figure *ast.CaptionFigure, entering bool) {}

func (r *Renderer) caption(w io.Writer, caption *ast.Caption, entering bool) {}

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
	case *mast.Bibliography:
	case *mast.BibliographyItem:
	case *mast.DocumentIndex, *mast.IndexLetter, *mast.IndexItem, *mast.IndexSubItem, *mast.IndexLink:
	case *ast.Footnotes:
		// do nothing, we're not outputing a footnote list
	case *ast.Text:
		r.text(w, node, entering)
	case *ast.Softbreak:
		// TODO
	case *ast.Hardbreak:
		r.hardBreak(w, node)
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
			r.outs(w, "\n.ti 0\n\\l'\\n(.lu'\n")
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
		r.outOneOf(w, true, "$", "$")
		if entering {
			r.out(w, node.Literal)
		}
		r.outOneOf(w, false, "$", "$")
	case *ast.Image:
		r.image(w, node, entering)
	case *ast.Code:
		r.outs(w, "\\fB\\fC")
		r.out(w, node.Literal)
		r.outs(w, "\\fR")
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

func (r *Renderer) callout(w io.Writer, node *ast.Callout, entering bool) {}

func (r *Renderer) text(w io.Writer, node *ast.Text, entering bool) {
	if !entering {
		return
	}
	text := node.Literal
	parent := node.Parent
	if parent != nil {
		if _, ok := parent.(*ast.Heading); ok {
			text = bytes.ToUpper(text)
		}
	}

	r.out(w, text)
}

func (r *Renderer) RenderHeader(w io.Writer, _ ast.Node) {
	r.outs(w, `.\" Generated by Mmark Markdown Processer - mmark.nl`+"\n")
}

func (r *Renderer) RenderFooter(w io.Writer, _ ast.Node) {}
