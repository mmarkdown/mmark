package xml3

// Renderer implements Renderer interface for IETF XMLv3 output.
type Renderer struct {
	opts RendererOptions
}

// New creates and configures an Renderer object, which satisfies the Renderer
// interface.
func New(opts RendererOptions) *Renderer {
	return &Renderer{opts: opts}
}

type RendererOptions struct {
	// Callouts are supported and detected by setting this option to the callout prefix.
	Callout string
}
