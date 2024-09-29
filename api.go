package progressy

import (
	"github.com/rollicks-c/progressy/internal/progbar"
	"io"
)

type ProgressBar interface {
	Complete()
	AddTask(name string) progbar.TaskHandler
}

func WithOverallTask(name string, stepCount int) progbar.Option {
	return progbar.WithOverallTask(name, stepCount)
}

func New(w io.Writer, options ...progbar.Option) ProgressBar {
	return progbar.New(w, options...)
}
