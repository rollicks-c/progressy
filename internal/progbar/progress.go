package progbar

import (
	"fmt"
	"sort"
	"strings"
)

type updater func(t task) task

func (p Progress) report() int {

	// task as ordinal list by line
	doneCount := 0
	taskList := make([]*task, 0, len(p.tasks))
	maxTaskWidth := 0
	for _, t := range p.tasks {
		taskList = append(taskList, t)
		if len(t.describe()) > maxTaskWidth {
			maxTaskWidth = len(t.describe())
		}
		if t.isComplete() {
			doneCount++
		}
	}
	if p.overallTask != nil {
		p.overallTask.current = doneCount
	}
	sort.Slice(taskList, func(i, j int) bool {
		return taskList[i].line < taskList[j].line
	})

	// print progress per task
	for _, t := range taskList {
		p.printProgress(*t, maxTaskWidth)
	}

	// report number of lines used
	return len(taskList)
}

func (p Progress) printProgress(task task, maxTaskWidth int) {

	// calc relative progress
	percentage := p.calcRelativeProgress(task.current, task.max)
	barWidth := 50
	progressWidth := (percentage * barWidth) / 100

	// clear line
	_, _ = fmt.Fprintf(p.w, "%s\r", strings.Repeat("", maxTaskWidth+barWidth+50))

	// print progress
	taskFill := strings.Repeat(" ", maxTaskWidth-len(task.describe()))
	taskName := task.describe()
	if task.isGlobal {
		taskName = fmt.Sprintf("\033[1m%s\033[0m", taskName)
	}
	p.printf("%s%s: ", taskName, taskFill)
	spaceFill := barWidth - progressWidth
	if spaceFill < 0 {
		spaceFill = 0
	}
	p.printf("|%s%s", strings.Repeat("#", progressWidth), strings.Repeat(" ", spaceFill))
	p.printf("| %d%%\n", percentage)

}

func (p Progress) moveCursorUp(lines int) {
	p.printf("\033[%dA", lines)
}

func (p Progress) printf(exp string, args ...interface{}) {
	_, _ = fmt.Fprintf(p.w, exp, args...)
}

func (p Progress) calcRelativeProgress(current, max int) int {

	if current > max {
		max = current
	}

	if max > 0 {
		return (current * 100) / max
	} else if current > max {
		return 100
	}
	return 0
}

func (p Progress) update(id taskID, h updater) int {

	p.lock.Lock()
	defer p.lock.Unlock()

	// update
	t := p.tasks[id]
	tNew := h(*t)
	tNew = p.sanitize(tNew)

	// preserve other props
	t.current = tNew.current
	t.max = tNew.max
	t.name = tNew.name
	p.tasks[id] = t

	lineCount := p.report()
	p.moveCursorUp(lineCount)

	return p.calcRelativeProgress(p.tasks[id].current, p.tasks[id].max)
}

func (p Progress) sanitize(t task) task {
	if t.current < 0 {
		t.current = 0
	}
	if t.max < 0 {
		t.max = 0
	}
	if t.current > t.max {
		t.current = t.max
	}
	return t
}

func (p Progress) complete(id taskID) {

	p.lock.Lock()
	defer p.lock.Unlock()

	p.tasks[id].current = p.tasks[id].max + 1

	lineCount := p.report()
	p.moveCursorUp(lineCount)
}
