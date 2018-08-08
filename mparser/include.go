// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted for mmark, by Miek Gieben, 2015.
// Adapted for mmark2 (fastly simplified and features removed), 2018.

package mparser

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// parseAddress parses a code address directive and returns the bytes or an error.
func parseAddress(addr []byte, data []byte) ([]byte, error) {
	bytes.TrimSpace(addr)

	if len(addr) == 0 {
		return data, nil
	}

	// check for prefix, either as ;prefix, prefix; or just standalone prefix.
	var prefix []byte
	if x := bytes.Index(addr, []byte("prefix=")); x >= 0 {
		if x+1 > len(addr) {
			return nil, fmt.Errorf("invalid prefix in address specification: %s", addr)
		}
		start := x + len("prefix=")
		quote := addr[start]
		if quote != '\'' && quote != '"' {
			return nil, fmt.Errorf("invalid prefix in address specification: %s", addr)
		}

		end := skipUntilChar(addr, start+1, quote)
		prefix = addr[start+1 : end]
		if len(prefix) == 0 {
			return nil, fmt.Errorf("invalid prefix in address specification: %s", addr)
		}

		addr = append(addr[:x], addr[end+1:]...)
		addr = bytes.Replace(addr, []byte(";"), []byte(""), 1)
		if len(addr) == 0 {
			data = addPrefix(data, prefix)
			return data, nil
		}
	}

	lo, hi, err := addrToByteRange(addr, data)
	if err != nil {
		return nil, err
	}

	// Acme pattern matches can stop mid-line,
	// so run to end of line in both directions if not at line start/end.
	for lo > 0 && data[lo-1] != '\n' {
		lo--
	}
	if hi > 0 {
		for hi < len(data) && data[hi-1] != '\n' {
			hi++
		}
	}

	data = data[lo:hi]
	if prefix != nil {
		data = addPrefix(data, prefix)
	}
	return data, nil
}

// addrToByteRange evaluates the given address. It returns the start and end index of the data we should return.
// Supported syntax:  N, M  or /start/, /end/ .
func addrToByteRange(addr, data []byte) (lo, hi int, err error) {
	chunk := bytes.Split(addr, []byte(","))
	if len(chunk) != 2 {
		return 0, 0, fmt.Errorf("invalid address specification: %s", addr)
	}
	left := bytes.TrimSpace(chunk[0])
	right := bytes.TrimSpace(chunk[1])

	if len(left) == 0 || len(right) == 0 {
		return 0, 0, fmt.Errorf("invalid address specification: %s", addr)
	}
	if left[0] == '/' { //regular expression
		if left[len(left)-1] != '/' {
			return 0, 0, fmt.Errorf("invalid address specification: %s", addr)
		}
		if right[0] != '/' {
			return 0, 0, fmt.Errorf("invalid address specification: %s", addr)
		}
		if right[len(right)-1] != '/' {
			return 0, 0, fmt.Errorf("invalid address specification: %s", addr)
		}

		lo, hi, err = addrRegexp(data, string(left[1:len(left)-1]), string(right[1:len(right)-1]))
		if err != nil {
			return 0, 0, err
		}
	} else {
		lo, err = strconv.Atoi(string(left))
		if err != nil {
			return 0, 0, err
		}
		i, j := 0, 0
		for i < len(data) {
			if data[i] == '\n' {
				j++
				if j >= lo {
					break
				}
			}
			i++
		}
		lo = i

		hi, err = strconv.Atoi(string(right))
		if err != nil {
			return 0, 0, err
		}
		i, j = 0, 0
		for i < len(data) {
			if data[i] == '\n' {
				j++
				if j+1 >= hi {
					break
				}
			}
			i++
		}
		hi = i
	}

	if lo > hi {
		return 0, 0, fmt.Errorf("invalid address specification: %s", addr)
	}

	return lo, hi, nil
}

// addrRegexp searches for pattern start and pattern end
func addrRegexp(data []byte, start, end string) (int, int, error) {
	start = "(?m:" + start + ")" // match through newlines
	reStart, err := regexp.Compile(start)
	if err != nil {
		return 0, 0, err
	}

	end = "(?m:" + end + ")"
	reEnd, err := regexp.Compile(end)
	if err != nil {
		return 0, 0, err
	}
	m := reStart.FindIndex(data)
	if len(m) == 0 {
		return 0, 0, errors.New("no match for " + start)
	}
	lo := m[0]

	m = reEnd.FindIndex(data[lo:]) // start *from* lo
	if len(m) == 0 {
		return 0, 0, errors.New("no match for " + end)
	}
	hi := m[0]

	return lo, hi, nil
}

func skipUntilChar(data []byte, i int, c byte) int {
	n := len(data)
	for i < n && data[i] != c {
		i++
	}
	return i
}

func addPrefix(data, prefix []byte) []byte {
	b := &bytes.Buffer{}
	b.Write(prefix)
	// assured that data ends in newline
	i := 0
	for i < len(data)-1 {
		b.WriteByte(data[i])
		if data[i] == '\n' {
			b.Write(prefix)
		}
		i++
	}
	return b.Bytes()
}
