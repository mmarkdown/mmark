package xml

import "bytes"

// Iscomment detects if a html span is a comment: <!--SPACE .... SPACE--!>. Return the comment text and true is so.
func IsComment(data []byte) ([]byte, bool) {
	if !bytes.HasPrefix(data, []byte("<!-- ")) {
		return nil, false
	}
	if !bytes.HasSuffix(data, []byte(" -->")) {
		return nil, false
	}

	return data[5 : len(data)-4], true
}
