package trace

import (
	"fmt"
	"io"
)

// Tracer describes an object capable of tracing events in code
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type nilTracer struct{}

func (t nilTracer) Trace(objects ...interface{}) {}

func (t tracer) Trace(objects ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(objects...)))
	t.out.Write([]byte("\n"))
}

// New creates a new default Tracer with the given Writer
func New(w io.Writer) Tracer {
	return &tracer{w}
}

// Off creates a tracer that does nothing with input
func Off() Tracer {
	return &nilTracer{}
}
