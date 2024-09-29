package progbar

import (
	"fmt"
	"io"
	"sync"
)

type TaskHandler struct {
	bar Progress
	id  taskID
}

type Progress struct {
	w     io.Writer
	tasks map[taskID]*task
	lock  *sync.Mutex

	overallTask *task
}

type Option func(p *Progress)

func WithOverallTask(name string, stepCount int) Option {
	return func(p *Progress) {
		u := p.AddTask(name)
		u.makeGlobal(p)
		u.Rel(0, stepCount)
	}
}

func New(w io.Writer, options ...Option) *Progress {
	p := &Progress{
		w:     w,
		tasks: make(map[taskID]*task),
		lock:  &sync.Mutex{},
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p Progress) AddTask(name string) TaskHandler {
	line := len(p.tasks)
	newTask := &task{
		line: line,
		name: name,
	}

	id := taskID(fmt.Sprintf("%d", line))
	p.tasks[id] = newTask

	updater := TaskHandler{
		bar: p,
		id:  id,
	}
	return updater
}

func (p Progress) Complete() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, t := range p.tasks {
		t.current = t.max + 1
	}

	p.report()
	p.printf("")
}

func (h TaskHandler) Complete() {
	h.bar.complete(h.id)
}

func (h TaskHandler) Abs(current, max int) int {
	return h.bar.update(h.id, func(t task) task {
		t.current = current
		t.max = max
		return t
	})
}

func (h TaskHandler) Rel(deltaCurrent, deltaMax int) int {
	return h.bar.update(h.id, func(t task) task {
		t.current += deltaCurrent
		t.max += deltaMax
		return t
	})
}

func (h TaskHandler) StepBy(step int) int {
	return h.bar.update(h.id, func(t task) task {
		t.current += step
		return t
	})
}

func (h TaskHandler) StepTo(value int) int {
	return h.bar.update(h.id, func(t task) task {
		t.current = value
		return t
	})
}
