package markdown

import "github.com/mmarkdown/mmark/mparser"

// Almost wholesale copy of parser/include.go - might make sense to make some of that public.

func isInclude(data []byte) bool {
	i := 0
	if data[i] != '{' || data[i+1] != '{' {
		return false
	}
	// find the end delimiter
	i = mparser.SkipUntilChar(data, i, '}')
	if i+1 >= len(data) {
		return false
	}
	end := i
	i++
	if data[i] != '}' {
		return false
	}

	if i+1 < len(data) && data[i+1] == '[' { // potential address specification
		start := i + 2
		end = mparser.SkipUntilChar(data, start, ']')
		if end >= len(data) {
			return false
		}
		return true
	}

	return true
}

func isCodeInclude(data []byte) bool {
	i := 0
	if len(data[i:]) < 3 {
		return false
	}
	if data[i] != '<' {
		return false
	}

	return isInclude(data[i+1:])
}
