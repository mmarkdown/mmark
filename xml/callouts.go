package xml

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// EscapeHTMLCallouts writes html-escaped d to w. It escapes &, <, > and " characters, *but*
// expands callouts <<N>> with the callout HTML, i.e. by calling rendering it as <N>.
func EscapeHTMLCallouts(w io.Writer, d []byte, comments [][]byte) {
	ld := len(d)
Parse:
	for i := 0; i < ld; i++ {
		for _, comment := range comments {
			if !bytes.HasPrefix(d[i:], comment) {
				break
			}

			lc := len(comment)
			if i+lc < ld {
				if id, consumed := parser.IsCallout(d[i+lc:]); consumed > 0 {
					// We have seen a callout
					io.WriteString(w, fmt.Sprintf("&lt;%s&gt;", id))
					i += consumed + lc - 1
					continue Parse
				}
			}
		}

		escSeq := html.Escaper[d[i]]
		if escSeq != nil {
			w.Write(escSeq)
		} else {
			w.Write([]byte{d[i]})
		}
	}
}
