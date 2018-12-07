package mast

// LatexSpan represents markdown LaTeX span node, i.e. any string that matches:
// \\[a-zA-Z]{.*}.
type LatexSpan struct {
	Leaf
}
