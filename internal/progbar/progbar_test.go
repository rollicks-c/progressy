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

}
