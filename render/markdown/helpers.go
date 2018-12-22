package markdown

import (
	"bytes"
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

func (r *Renderer) push(data []byte) { r.prefix.push(data) }
func (r *Renderer) peek() int        { return r.prefix.peek() }

func (p *prefixStack) push(data []byte) { p.p = append(p.p, data) }

func (p *prefixStack) pop() []byte {
	if len(p.p) == 0 {
		return nil
	}
	last := p.p[len(p.p)-1]
	p.p = p.p[:len(p.p)-1]
	return last
}

// peek returns the lenght of the last pushed element.
func (p *prefixStack) peek() int {
	if len(p.p) == 0 {
		return 0
	}
	last := p.p[len(p.p)-1]
	return len(last)
}

// flatten flattens the stack in reverse order.
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

// listPrefixLength returns the length of the prefix we need for list in ast.Node
func listPrefixLength(list *ast.List, start int) int {
	numChild := len(list.Children) + start
	switch {
	case numChild < 10:
		return 3
	case numChild < 100:
		return 4
	case numChild < 1000:
		return 5
	}
	return 6 // bit of a ridicules list
}

func Space(length int) []byte { return bytes.Repeat([]byte(" "), length) }

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}
