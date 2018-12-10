package mparser

import "testing"

func TestLaTeXHook(t *testing.T) {
	data := []byte(`\citec{citation}`)
	ast, _, _ := LatexHook(data)
	if ast == nil {
		t.Errorf("Expected ast.LatexSpan, got none on onput: %s", data)
	}
	if string(ast.AsLeaf().Content) != string(data) {
		t.Errorf("Expected ast.LatexSpan, got none on onput: %s", data)
	}

	data = []byte(`\citec0{citation}`)
	ast, _, _ = LatexHook(data)
	if ast == nil {
		t.Errorf("Expected ast.LatexSpan, got none on onput: %s", data)
	}
	if string(ast.AsLeaf().Content) != string(data) {
		t.Errorf("Expected ast.LatexSpan, got none on onput: %s", data)
	}

	data = []byte(`\{citation}`)
	ast, _, _ = LatexHook(data)
	if ast != nil {
		t.Errorf("Got ast on invalid input: %s", data)
	}
}
