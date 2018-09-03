package xml

import "bytes"

// words2119 contains the words we should recognize as BCP 14 words when used with **strong**.
var words2119 = [][]byte{
	[]byte("MUST"),
	[]byte("MUST NOT"),
	[]byte("REQUIRED"),
	[]byte("SHALL"),
	[]byte("SHALL NOT"),
	[]byte("SHOULD"),
	[]byte("SHOULD NOT"),
	[]byte("RECOMMENDED"),
	[]byte("NOT RECOMMENDED"),
	[]byte("MAY"),
	[]byte("OPTIONAL"),
}

// Is2119 checks if word is a RFC 2119 word.
func Is2119(word []byte) bool {
	for _, bcp := range words2119 {
		if bytes.Compare(word, bcp) == 0 {
			return true
		}
	}
	return false
}
