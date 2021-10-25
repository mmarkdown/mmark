package main

// run throught the first few bytes of d and check to see if the line
// ending contains \n  or \cr \n  - if the latter is found return true.
func crlf(d []byte) bool {
	found := false
	for i := range d {
		if d[i] == '\r' {
			found = true
		}
		if d[i] == '\n' {
			return found
		}
	}

	return false
}
