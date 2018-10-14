package markdown

import "github.com/mmarkdown/mmark/mparser"

// Almost wholesale copy of parser/include.go - might make sense to make some of that public.

func isInclude(data []byte) int {
	i := 0
	if data[i] != '{' || data[i+1] != '{' {
		return 0
	}
	// find the end delimiter
	i = mparser.SkipUntilChar(data, i, '}')
	if i+1 >= len(data) {
		return 0
	}
	end := i
	i++
	if data[i] != '}' {
		return 0
	}

	if i+1 < len(data) && data[i+1] == '[' { // potential address specification
		start := i + 2
		end = mparser.SkipUntilChar(data, start, ']')
		if end >= len(data) {
			return 0
		}
		return end
	}

	return i
}

func isCodeInclude(data []byte) int {
	i := 0
	if len(data[i:]) < 3 {
		return 0
	}
	if data[i] != '<' {
		return 0
	}

	x := isInclude(data[i+1:])
	if x == 0 {
		return 0
	}
	return x + 2
}
