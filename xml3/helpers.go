package xml3

import (
	"fmt"
	"io"
	"strings"
)

func (r *Renderer) out(w io.Writer, d []byte) {
	w.Write(d)
}

func (r *Renderer) outs(w io.Writer, s string) {
	io.WriteString(w, s)
}

func (r *Renderer) cr(w io.Writer) {
	r.outs(w, "\n")
}

func (r *Renderer) outTag(w io.Writer, name string, attrs []string) {
	s := name
	if len(attrs) > 0 {
		s += " " + strings.Join(attrs, " ")
	}
	io.WriteString(w, s+">")
}

// outTagContents output the opening tag with possible attributes, then the content
// and then the closing tag.
func (r *Renderer) outTagContent(w io.Writer, name string, attrs []string, content string) {
	s := name
	if len(attrs) > 0 {
		s += " " + strings.Join(attrs, " ")
	}
	io.WriteString(w, s+">")
	io.WriteString(w, content)
	io.WriteString(w, "</"+name[1:]+">\n")
}

func attributes(keys, values []string) (s []string) {
	for i, k := range keys {
		if values[i] == "" { // skip entire k=v is value is empty
			continue
		}
		s = append(s, fmt.Sprintf(`%s="%s"`, k, values[i]))
	}
	return s
}
