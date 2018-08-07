package xml

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/html"
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

func (r *Renderer) outOneOf(w io.Writer, outFirst bool, first string, second string) {
	if outFirst {
		r.outs(w, first)
	} else {
		r.outs(w, second)
	}
}

func (r *Renderer) outOneOfCr(w io.Writer, outFirst bool, first string, second string) {
	if outFirst {
		r.cr(w)
		r.outs(w, first)
	} else {
		r.outs(w, second)
		r.cr(w)
	}
}

// outTagContents output the opening tag with possible attributes, then the content
// and then the closing tag.
func (r *Renderer) outTagContent(w io.Writer, name string, attrs []string, content string) {
	s := name
	if len(attrs) > 0 {
		s += " " + strings.Join(attrs, " ")
	}
	io.WriteString(w, s+">")
	html.EscapeHTML(w, []byte(content))
	io.WriteString(w, "</"+name[1:]+">\n")
}

func (r *Renderer) sectionClose(w io.Writer) {
	if r.section == nil {
		return
	}

	tag := "</section>"
	if r.section.Special != nil {
		tag = "</note>"
		if isAbstract(r.section.Special) {
			tag = "</abstract>"
		}
	}
	r.outs(w, tag)
	r.cr(w)
}

func attributes(keys, values []string) (s []string) {
	for i, k := range keys {
		if values[i] == "" { // skip entire k=v is value is empty
			continue
		}
		v := escapeHTMLString(values[i])
		s = append(s, fmt.Sprintf(`%s="%s"`, k, v))
	}
	return s
}

func isAbstract(word []byte) bool {
	return strings.EqualFold(string(word), "abstract")
}

func escapeHTMLString(s string) string {
	buf := &bytes.Buffer{}
	html.EscapeHTML(buf, []byte(s))
	return buf.String()
}

func appendLanguageAttr(attrs []string, info []byte) []string {
	if len(info) == 0 {
		return attrs
	}
	endOfLang := bytes.IndexAny(info, "\t ")
	if endOfLang < 0 {
		endOfLang = len(info)
	}
	s := `type="` + string(info[:endOfLang]) + `"`
	return append(attrs, s)
}
