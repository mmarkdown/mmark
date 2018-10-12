package markdown

import (
	"bytes"
	"io"
	"regexp"
	"unicode"

	"github.com/gomarkdown/markdown/ast"
	"github.com/kr/text"
)

func (r *Renderer) outOneOf(w io.Writer, outFirst bool, first, second string) {
	if outFirst {
		r.outs(w, first)
	} else {
		r.outs(w, second)
	}
}

func (r *Renderer) out(w io.Writer, d []byte)  { w.Write(d) }
func (r *Renderer) outs(w io.Writer, s string) { io.WriteString(w, s) }
func (r *Renderer) cr(w io.Writer)             { r.outs(w, "\n") }
func (r *Renderer) outPrefix(w io.Writer)      { r.out(w, r.prefix.flatten()) }

func (r *Renderer) newline(w io.Writer) {
	r.out(w, r.prefix.flatten())
	r.outs(w, "\n")
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

// escapeText escape the text in data using isEscape.
func escapeText(data []byte) []byte {
	buf := &bytes.Buffer{}

	for i := range data {
		switch data[i] {
		case '<':
			if isCodeInclude(data[i:]) {
				buf.WriteByte(data[i])
				continue
			}
			fallthrough
		case '>':
			fallthrough
		case '&':
			fallthrough
		case '\\':
			if !isEscape(data, i) {
				buf.WriteByte('\\')
			}
		}
		buf.WriteByte(data[i])
	}
	return buf.Bytes()
}

// isEscape returns true if byte i is prefixed by an odd number of backslahses.
func isEscape(data []byte, i int) bool {
	if i == 0 {
		return false
	}
	if i == 1 {
		return data[0] == '\\'
	}
	j := i - 1
	for ; j >= 0; j-- {
		if data[j] != '\\' {
			break
		}
	}
	j++
	// odd number of backslahes means escape
	return (i-j)%2 != 0
}

// Copied from gomarkdown/markdown.

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
