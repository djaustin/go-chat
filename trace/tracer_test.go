package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Got nil from New, expected Tracer")
	} else {
		tracer.Trace("trace test")
		if buf.String() != "trace test\n" {
			t.Errorf("Got '%s' from Trace, expected 'trace test\n'", buf.String())
		}
	}
}
