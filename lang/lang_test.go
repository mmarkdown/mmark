package lang

import "testing"

func TestBibliography(t *testing.T) {
	l := New("en")
	if l.Bibliography() != "Bibliography" {
		t.Errorf("expected %s, got %s", "Bibliography", l.Bibliography())
	}
	l = New("not-defined")
	if l.Bibliography() != "Bibliography" {
		t.Errorf("expected %s, got %s", "Bibliography", l.Bibliography())
	}
}
