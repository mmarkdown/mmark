package xml3

// Renderer implements Renderer interface for IETF XMLv3 output.
type Renderer struct{}

// New creates and configures an Renderer object, which satisfies the Renderer
// interface.
func New() *Renderer { return &Renderer{} }
