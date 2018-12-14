package latex

import (
	"io"
	"strings"
)

var special = map[byte]struct{}{
	'_':  struct{}{},
	'{':  struct{}{},
	'}':  struct{}{},
	'%':  struct{}{},
	'$':  struct{}{},
	'&':  struct{}{},
	'\\': struct{}{},
	'~':  struct{}{},
	'#':  struct{}{},
}

func escapeSpecialChars(out io.Writer, text []byte) {
	for i := 0; i < len(text); i++ {
		if _, isSpc := special[text[i]]; isSpc {
			out.Write([]byte("\\"))
		}
		out.Write([]byte{text[i]})
	}
}

// IsAbstract returns if word is equal to abstract.
func IsAbstract(word []byte) bool              { return strings.EqualFold(string(word), "abstract") }
func (r *Renderer) out(w io.Writer, d []byte)  { w.Write(d) }
func (r *Renderer) outs(w io.Writer, s string) { io.WriteString(w, s) }
func (r *Renderer) cr(w io.Writer)             { r.outs(w, "\n") }

func (r *Renderer) outOneOf(w io.Writer, outFirst bool, first string, second string) {
	if outFirst {
		r.cr(w)
		r.outs(w, first)
	} else {
		r.outs(w, second)
	}
}
