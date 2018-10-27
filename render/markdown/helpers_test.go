package markdown

import (
	"bytes"
	"testing"
)

func TestPrefixPush(t *testing.T) {
	prefix := &prefixStack{p: [][]byte{}}
	prefix.push([]byte("A"))
	prefix.push([]byte("B"))

	if bytes.Compare(prefix.flatten(), []byte("AB")) != 0 {
		t.Errorf("Expected %s, got %s", "AB", prefix.flatten())
	}

	prefix.pop()

	if bytes.Compare(prefix.flatten(), []byte("A")) != 0 {
		t.Errorf("Expected %s, got %s", "A", prefix.flatten())
	}
}
