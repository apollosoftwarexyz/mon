package mon

import (
	"sync/atomic"
	"time"
)

// TaskBuilder for adding new [Task] references to a monitor ([M]).
type TaskBuilder interface {
	// Name sets the name of the task.
	Name(name string) TaskBuilder

	// Caption sets the caption of the task.
	Caption(caption string) TaskBuilder

	// Category sets the category of the task.
	Category(category string) TaskBuilder

	// TotalSteps sets the total number of steps that must be completed as part
	// of this task. If this is not set, then [Task.IsIndeterminate] is true.
	TotalSteps(totalSteps uint64) TaskBuilder

	// Apply the task to the monitor that created the builder.
	//
	// This is the terminal step of the builder and returns the [Task] reference
	// that was added to the monitor.
	Apply() Task
}

type taskBuilder struct {
	m          *model
	name       string
	caption    string
	category   string
	totalSteps uint64
}

func (b *taskBuilder) Name(name string) TaskBuilder {
	b.name = name
	return b
}

func (b *taskBuilder) Caption(caption string) TaskBuilder {
	b.caption = caption
	return b
}

func (b *taskBuilder) Category(category string) TaskBuilder {
	b.category = category
	return b
}

func (b *taskBuilder) TotalSteps(totalSteps uint64) TaskBuilder {
	b.totalSteps = totalSteps
	return b
}

func (b *taskBuilder) Apply() Task {
	stepsTotal := &atomic.Uint64{}
	stepsTotal.Store(b.totalSteps)

	task := &task{
		notify:         b.m.notify,
		name:           b.name,
		caption:        b.caption,
		category:       b.category,
		startTime:      time.Now(),
		stepsCompleted: &atomic.Uint64{},
		stepsTotal:     stepsTotal,
	}
	b.m.addTask(task)
	return task
}

// Task tracked by a monitor, [M].
type Task interface {
	// GetName of the task.
	GetName() string

	// GetCaption of the task.
	GetCaption() string

	// GetCategory of the task.
	GetCategory() string

	// GetStartedAt returns the time that the task was started at.
	GetStartedAt() time.Time

	// GetCompletedAt returns the time that the task was completed at. If the
	// task is not completed (i.e., IsCompleted is false), this function returns
	// a time for which [time.Time.IsZero] returns true.
	GetCompletedAt() time.Time

	// GetElapsed time since the task was created and before the task was
	// stopped.
	GetElapsed() time.Duration

	// GetProgress expressed as a percentage. For tasks where IsIndeterminate is
	// true, this is always zero.
	GetProgress() float64

	// GetAverageTimePerStep computes a mean average of time per step using up
	// to 256 discrete previous step times. If there are no completed steps,
	// this function returns zero and false.
	GetAverageTimePerStep() (time.Duration, bool)

	// GetEstimatedCompletion duration from now.
	//
	// This function uses [Task.GetAverageTimePerStep] and the number of
	// remaining steps to extrapolate a completion time. This relies on an
	// assumption that steps are sequential within a task and that steps are
	// equal (that is, they should nominally take roughly the same amount of
	// time to complete).
	//
	// If there are no completed steps, or the task is already complete, this
	// function returns zero and false.
	GetEstimatedCompletion() (time.Duration, bool)

	// IsIndeterminate indicates whether a total number of steps is known.
	//
	// If it is not (i.e., StepsTotal is 0) then the task is indeterminate.
	IsIndeterminate() bool

	// IsCompleted indicates that a task is fully completed. If the task
	// IsIndeterminate, this is true if there are any completed steps.
	IsCompleted() bool

	// CompletedStep increments the number of steps that have already been
	// completed as part of this task.
	//
	// If the task IsDone, this function is a no-op.
	CompletedStep()

	// GetCompletedSteps returns the number of steps that have already been
	// completed as part of this task.
	GetCompletedSteps() uint64

	// CompletedSteps sets the number of steps that have already been completed
	// as part of this task.
	//
	// If the task IsDone, this function is a no-op.
	CompletedSteps(completedSteps uint64)

	// GetTotalSteps returns the number of steps that must be completed as part
	// of this task.
	GetTotalSteps() uint64

	// TotalSteps sets the total number of steps that must be completed as part
	// of this task.
	//
	// If the task IsDone, this function is a no-op.
	TotalSteps(totalSteps uint64)
}

type notifyFn func()

type task struct {
	notify         notifyFn
	name           string
	caption        string
	category       string
	startTime      time.Time
	endTime        time.Time
	stepsCompleted *atomic.Uint64
	stepsTotal     *atomic.Uint64

	timeOfLastRecord time.Time
	timePerStep      []time.Duration
}

func (t *task) GetName() string           { return t.name }
func (t *task) GetCaption() string        { return t.caption }
func (t *task) GetCategory() string       { return t.category }
func (t *task) GetStartedAt() time.Time   { return t.startTime }
func (t *task) GetCompletedAt() time.Time { return t.endTime }

func (t *task) GetElapsed() time.Duration {
	if !t.endTime.IsZero() {
		return t.endTime.Sub(t.startTime)
	}

	return time.Since(t.startTime)
}

func (t *task) GetProgress() float64 {
	total := t.stepsTotal.Load()
	if total == 0 {
		return 0
	}

	completed := t.stepsCompleted.Load()
	return float64(completed) / float64(total)
}

func (t *task) GetAverageTimePerStep() (time.Duration, bool) {
	if len(t.timePerStep) == 0 {
		return 0, false
	}

	var avgTimePerStep time.Duration
	for _, metric := range t.timePerStep {
		avgTimePerStep += metric
	}
	avgTimePerStep /= time.Duration(len(t.timePerStep))
	return avgTimePerStep, true
}

func (t *task) GetEstimatedCompletion() (time.Duration, bool) {
	avgTimePerStep, ok := t.GetAverageTimePerStep()

	if !ok || t.IsCompleted() {
		return 0, false
	}

	remainingSteps := t.GetTotalSteps() - t.GetCompletedSteps()
	return time.Duration(remainingSteps) * avgTimePerStep, true
}

func (t *task) IsIndeterminate() bool { return t.stepsTotal.Load() == 0 }
func (t *task) IsCompleted() bool     { return !t.endTime.IsZero() }

func (t *task) recordTimePerSteps(n uint64) {
	var d time.Duration
	if t.timeOfLastRecord.IsZero() {
		d = time.Since(t.startTime)
	} else {
		d = time.Since(t.timeOfLastRecord)
	}

	if n < 1 {
		return
	}

	if n > 1 {
		d /= time.Duration(n)
	}

	if len(t.timePerStep) >= 256 {
		t.timeOfLastRecord = time.Now()
		t.timePerStep = append(t.timePerStep[1:], d)
	} else {
		t.timeOfLastRecord = time.Now()
		t.timePerStep = append(t.timePerStep, d)
	}
}

func (t *task) checkCompleted() {
	completed := t.stepsCompleted.Load()
	total := t.stepsTotal.Load()

	var isDone bool
	if total == 0 {
		isDone = completed > 0
	} else {
		isDone = completed >= total
	}

	if isDone {
		t.endTime = time.Now()
	}
}

func (t *task) CompletedStep() {
	if t.IsCompleted() {
		return
	}

	t.stepsCompleted.Add(1)
	t.recordTimePerSteps(1)
	t.checkCompleted()
	t.notify()
}

func (t *task) GetCompletedSteps() uint64 {
	return t.stepsCompleted.Load()
}

func (t *task) CompletedSteps(completedSteps uint64) {
	if t.IsCompleted() {
		return
	}

	previouslyCompletedSteps := t.stepsCompleted.Swap(completedSteps)
	t.recordTimePerSteps(previouslyCompletedSteps)
	t.checkCompleted()
	t.notify()
}

func (t *task) GetTotalSteps() uint64 {
	return t.stepsTotal.Load()
}

func (t *task) TotalSteps(totalSteps uint64) {
	if t.IsCompleted() {
		return
	}

	t.stepsTotal.Store(totalSteps)
	t.checkCompleted()
	t.notify()
}
