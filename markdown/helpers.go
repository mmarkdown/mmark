package markdown

import (
	"io"
	"regexp"
	"unicode"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/internal/text"
)

func (r *Renderer) outOneOf(w io.Writer, outFirst bool, first, second string) {
	if outFirst {
		r.outs(w, first)
	} else {
		r.outs(w, second)
	}
}

func (r *Renderer) out(w io.Writer, d []byte)  { w.Write(d); r.suppress = false }
func (r *Renderer) outs(w io.Writer, s string) { io.WriteString(w, s); r.suppress = false }
func (r *Renderer) outPrefix(w io.Writer)      { r.out(w, r.prefix.flatten()); r.suppress = false }
func (r *Renderer) endline(w io.Writer)        { r.outs(w, "\n"); r.suppress = false }

func (r *Renderer) newline(w io.Writer) {
	if r.suppress {
		return
	}
	r.out(w, r.prefix.flatten())
	r.outs(w, "\n")
	r.suppress = true
}

var re = regexp.MustCompile("  +")

// lastNode returns true if we are the last node under this parent.
func lastNode(node ast.Node) bool { return ast.GetNextNode(node) == nil }

// wrapText wraps the text in data, taking len(prefix) into account.
func (r *Renderer) wrapText(data, prefix []byte) []byte {
	replaced := re.ReplaceAll(data, []byte(" "))
	wrapped := text.WrapBytes(replaced, r.opts.TextWidth-len(prefix))
	return r.indentText(wrapped, prefix)
}

func (r *Renderer) indentText(data, prefix []byte) []byte {
	return text.IndentBytes(data, prefix)
}

// sanitizeAnchorName returns a sanitized anchor name for the given text.
// Taken from https://github.com/shurcooL/sanitized_anchor_name/blob/master/main.go#L14:1
func sanitizeAnchorName(text string) string {
	var anchorName []rune
	var futureDash = false
	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
}

type prefixStack struct {
	p [][]byte
}

func (r *Renderer) pop() []byte {
	last := r.prefix.pop()
	if last != nil && r.prefix.len() == 0 {
		r.suppress = false
	}
	return last
}

func (r *Renderer) push(data []byte) {
	r.prefix.push(data)
}

func (p *prefixStack) push(data []byte) { p.p = append(p.p, data) }

func (p *prefixStack) pop() []byte {
	if len(p.p) == 0 {
		return nil
	}
	last := p.p[len(p.p)-1]
	p.p = p.p[:len(p.p)-1]
	return last
}

// flatten stack in reverse order
func (p *prefixStack) flatten() []byte {
	ret := []byte{}
	for _, b := range p.p {
		ret = append(ret, b...)
	}
	return ret
}

func (p *prefixStack) len() (l int) {
	for _, b := range p.p {
		l += len(b)
	}
	return l
}
