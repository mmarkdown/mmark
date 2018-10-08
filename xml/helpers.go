package xml

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/mmarkdown/mmark/mast"
)

func (r *Renderer) out(w io.Writer, d []byte)  { w.Write(d) }
func (r *Renderer) outs(w io.Writer, s string) { io.WriteString(w, s) }
func (r *Renderer) cr(w io.Writer)             { r.outs(w, "\n") }

func (r *Renderer) outAttr(w io.Writer, attrs []string) {
	if len(attrs) > 0 {
		io.WriteString(w, " ")
		io.WriteString(w, strings.Join(attrs, " "))
	}
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

func (r *Renderer) outTagMaybe(w io.Writer, name string, content string) {
	if content != "" {
		r.outTagContent(w, name, content)
	}
}

func (r *Renderer) outTagContent(w io.Writer, name string, content string) {
	io.WriteString(w, name+">")
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

	// TODO(miek): Probably better to actually count the number of OPEN sections instead of the level of them.
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

func appendLanguageAttr(node ast.Node, info []byte) {
	if len(info) == 0 {
		return
	}
	endOfLang := bytes.IndexAny(info, "\t ")
	if endOfLang < 0 {
		endOfLang = len(info)
	}
	mast.SetAttribute(node, "type", info[:endOfLang])
}

// Attributes returns the key values in a stringslice where each is type set as key="value".
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

// AttributesContains checks if the attribute list contains key.
func AttributesContains(key string, attr []string) bool {
	check := key + `="`
	for _, a := range attr {
		if strings.HasPrefix(a, check) {
			return true
		}
	}
	return false
}
