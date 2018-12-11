package latex

import (
	"bytes"
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

func escapeSpecialChars(out *bytes.Buffer, text []byte) {
	for i := 0; i < len(text); i++ {
		// directly copy normal characters
		org := i

		_, isSpc := special[text[i]]
		for i < len(text) && !isSpc {
			i++
		}
		if i > org {
			out.Write(text[org:i])
		}

		// escape a character
		if i >= len(text) {
			break
		}
		out.WriteByte('\\')
		out.WriteByte(text[i])
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
