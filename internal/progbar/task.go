package progbar

import "fmt"

type taskID string

type task struct {
	line    int
	name    string
	current int
	max     int

	isGlobal bool
}

func (t task) describe() string {
	return fmt.Sprintf("%s (%d)", t.name, t.max)
}

func (t task) isComplete() bool {
	return t.current >= t.max
}

func (h TaskHandler) makeGlobal(p *Progress) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.overallTask = p.tasks[h.id]
	p.overallTask.isGlobal = true
}
