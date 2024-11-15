package text

import (
	"bytes"

	mtext "github.com/mmarkdown/mmark/v2/internal/text"
)

func (r *Renderer) centerText(data []byte) []byte {
	replaced := re.ReplaceAll(data, []byte(" "))
	wrapped := mtext.WrapBytes(replaced, []byte{}, r.opts.TextWidth)

	// now split the wrapped text on the newlines and center each chucnk

	p := bytes.Split(wrapped, []byte("\n"))
	var centered []byte
	for i := range p {
		left := r.opts.TextWidth - len(p[i])
		if left <= 0 {
			centered = append(centered, p[i]...)
			centered = append(centered, '\n')
			continue
		}
		centered = append(centered, bytes.Repeat([]byte{' '}, left/2)...)
		centered = append(centered, p[i]...)
		centered = append(centered, '\n')
	}
	return centered
}
