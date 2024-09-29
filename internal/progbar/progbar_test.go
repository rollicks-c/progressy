package progbar

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestSanitize(t *testing.T) {
	h := New(io.Discard).AddTask("")

	// normal
	assert.Equal(t, 0, h.Abs(0, 10))
	assert.Equal(t, 50, h.Abs(5, 10))
	assert.Equal(t, 100, h.Abs(10, 10))

	// invalid
	assert.Equal(t, 100, h.Abs(11, 10))
	assert.Equal(t, 0, h.Abs(-1, 10))
	assert.Equal(t, 0, h.Abs(0, -1))
	assert.Equal(t, 0, h.Abs(-1, -1))

}

func TestReport(t *testing.T) {

	out := bytes.NewBufferString("")
	h := New(out).AddTask("task one")

	h.Abs(0, 10)
	assert.Equal(t, "\rtask one (10): |                                                  | 0%\n\u001B[1A", out.String())

	out.Reset()
	h.Abs(1, 10)
	assert.Equal(t, "\rtask one (10): |#####                                             | 10%\n\u001B[1A", out.String())

	out.Reset()
	h.Abs(5, 10)
	assert.Equal(t, "\rtask one (10): |#########################                         | 50%\n\u001B[1A", out.String())

	out.Reset()
	h.Abs(10, 10)
	assert.Equal(t, "\rtask one (10): |##################################################| 100%\n\u001B[1A", out.String())

	out.Reset()
	h.Abs(11, 10)
	assert.Equal(t, "\rtask one (10): |##################################################| 100%\n\u001B[1A", out.String())

	h.Abs(15, 15)
	out.Reset()
	h.Complete()
	assert.Equal(t, "\rtask one (15): |##################################################| 100%\n\u001B[1A", out.String())

	pg := New(out)
	pg.AddTask("task two")
	out.Reset()
	pg.Complete()
	assert.Equal(t, "\rtask two (0): |##################################################| 100%\n", out.String())

}

func TestProgress(t *testing.T) {
	h := New(io.Discard).AddTask("")

	// rel
	assert.Equal(t, 0, h.Abs(0, 10))
	assert.Equal(t, 25, h.Rel(5, 10))
	assert.Equal(t, 50, h.Rel(5, 0))

	// step rel
	assert.Equal(t, 0, h.Abs(0, 10))
	assert.Equal(t, 0, h.StepBy(0))
	assert.Equal(t, 50, h.StepBy(5))
	assert.Equal(t, 100, h.StepBy(10))
	assert.Equal(t, 100, h.StepBy(10))

	// step abs
	assert.Equal(t, 0, h.Abs(0, 10))
	assert.Equal(t, 0, h.StepTo(0))
	assert.Equal(t, 50, h.StepTo(5))
	assert.Equal(t, 100, h.StepTo(10))
	assert.Equal(t, 100, h.StepTo(20))

}

func TestGlobalTask(t *testing.T) {

	out := bytes.NewBufferString("")
	h := New(out, WithOverallTask("global", 15)).AddTask("task one")
	assert.Equal(t, "\r\x1b[1mglobal (15)\x1b[0m: |                                                  | 0%\n\x1b[1A", out.String())

	out.Reset()
	h.Abs(0, 10)
	assert.Equal(t, "\r\x1b[1mglobal (15)\x1b[0m  : |                                                  | 0%\n\rtask one (10): |                                                  | 0%\n\x1b[2A", out.String())

}
