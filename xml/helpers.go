package xml

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

func (r *Renderer) out(w io.Writer, d []byte)  { w.Write(d) }
func (r *Renderer) outs(w io.Writer, s string) { io.WriteString(w, s) }
func (r *Renderer) cr(w io.Writer)             { r.outs(w, "\n") }

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

func (r *Renderer) sectionClose(w io.Writer, new *ast.Heading) {
	defer func() {
		r.section = new
	}()

	if r.section == nil {
		return
	}

	if r.section.IsSpecial {
		tag := "</note>"
		if IsAbstract(r.section.Literal) {
			tag = "</abstract>"
		}
		r.outs(w, tag)
		r.cr(w)
		return
	}

	tag := "</section>"
	curLevel := r.section.Level
	newLevel := 1 // close them all
	if new != nil {
		newLevel = new.Level
	}

	// subheading in a heading
	if newLevel > curLevel {
		return
	}

	for i := newLevel; i <= curLevel; i++ {
		r.outs(w, tag)
		r.cr(w)
	}
}

func (r *Renderer) ensureUniqueHeadingID(id string) string {
	for count, found := r.headingIDs[id]; found; count, found = r.headingIDs[id] {
		tmp := fmt.Sprintf("%s-%d", id, count+1)

		if _, tmpFound := r.headingIDs[tmp]; !tmpFound {
			r.headingIDs[id] = count + 1
			id = tmp
		} else {
			id = id + "-1"
		}
	}

	if _, found := r.headingIDs[id]; !found {
		r.headingIDs[id] = 0
	}

	return id
}

// Attributes returns the key values in a stringslice where key="value".
func Attributes(keys, values []string) (s []string) {
	for i, k := range keys {
		if values[i] == "" { // skip entire k=v is value is empty
			continue
		}
		v := EscapeHTMLString(values[i])
		s = append(s, fmt.Sprintf(`%s="%s"`, k, v))
	}
	return s
}

// IsAbstract returns if word is equal to abstract.
func IsAbstract(word []byte) bool {
	return strings.EqualFold(string(word), "abstract")
}

// EscapeHTMLString escapes the string s.
func EscapeHTMLString(s string) string {
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
