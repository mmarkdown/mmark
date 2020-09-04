package xml

import "bytes"

// Iscomment detects if a html span is a comment: <!-- .... --!>. Return the comment text and true is so.
func IsComment(data []byte) ([]byte, bool) {
	if !bytes.HasPrefix(data, []byte("<!--")) {
		return nil, false
	}
	if !bytes.HasSuffix(data, []byte("-->")) {
		return nil, false
	}

	return data[5 : len(data)-4], true
}

func IsBr(data []byte) bool {
	// <br> <br/> <br /> and <br></br> are recognized
	if bytes.Equal(data, []byte("<br>")) {
		return true
	}
	if bytes.Equal(data, []byte("<br >")) {
		return true
	}
	if bytes.Equal(data, []byte("<br/>")) {
		return true
	}
	if bytes.Equal(data, []byte("<br />")) {
		return true
	}
	if bytes.Equal(data, []byte("<br></br>")) {
		return true
	}
	return false
}
