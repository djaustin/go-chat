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

func (t tracer) Trace(objects ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(objects...)))
	t.out.Write([]byte("\n"))
}

// New creates a new default Tracer with the given Writer
func New(w io.Writer) Tracer {
	return &tracer{w}
}
