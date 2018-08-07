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

	return data[lo:hi], nil
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
		if right[0] != '/' {
			return 0, 0, fmt.Errorf("invalid address specification: %s", addr)
		}

		lo, hi, err = addrRegexp(data, string(left), string(right))
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
				if j >= hi {
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
	reStart, err := regexp.Compile(start)
	if err != nil {
		return 0, 0, err
	}
	reEnd, err := regexp.Compile(end)
	if err != nil {
		return 0, 0, err
	}
	m := reStart.FindIndex(data)
	if len(m) == 0 {
		return 0, 0, errors.New("no match for " + start)
	}
	lo := m[0]

	m = reEnd.FindIndex(data)
	if len(m) == 0 {
		return 0, 0, errors.New("no match for " + end)
	}
	hi := m[0]

	return lo, hi, nil
}
