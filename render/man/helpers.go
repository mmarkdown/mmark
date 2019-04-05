package man

import (
	"io"
)

func (r *Renderer) out(w io.Writer, d []byte)  { w.Write(d) }
func (r *Renderer) outs(w io.Writer, s string) { io.WriteString(w, s) }

func (r *Renderer) outOneOf(w io.Writer, outFirst bool, first string, second string) {
	if outFirst {
		r.outs(w, first)
	} else {
		r.outs(w, second)
	}
}

func needsBackslash(c byte) bool {
	for _, r := range []byte("-_&\\") {
		if c == r {
			return true
		}
	}
	return false
}

func escapeSpecialChars(r *Renderer, w io.Writer, text []byte) {
	for i := 0; i < len(text); i++ {
		// escape apostrophe or period after newline (making this first char on the line)
		if i == 0 && (text[i] == '\'' || text[i] == '.') {
			r.outs(w, "\\&")
			r.out(w, []byte{text[i]})
			continue
		}

		if i > 0 && text[i-1] == '\n' && (text[i] == '\'' || text[i] == '.') {
			r.outs(w, "\\&")
			r.out(w, []byte{text[i]})
			continue
		}

		if needsBackslash(text[i]) {
			r.out(w, []byte{'\\'})
		}
		r.out(w, []byte{text[i]})
	}
}
