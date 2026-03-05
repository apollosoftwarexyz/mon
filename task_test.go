package mon_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/apollosoftwarexyz/mon"
	"github.com/apollosoftwarexyz/mon/formatting"
	"github.com/stretchr/testify/assert"
)

const (
	notMockName    = "not name"
	mockName       = "name"
	mockCaption    = "caption"
	mockCategory   = "category"
	mockTotalSteps = uint64(100)
)

var (
	mockUnit  = &formatting.BytesUnit{}
	mockError = fmt.Errorf("mock error")
)

// createDefaultTask creates a new task on a new monitor with empty values.
func createDefaultTask() mon.Task {
	m := mon.New("test")
	return m.AddTask().Apply()
}

// TestAddTask performs a smoke test of the task builder to ensure it assigns
// values correctly.
func TestAddTask(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().
		// test that builder methods overwrite former build methods
		Name(notMockName).
		Name(mockName).
		Caption(mockCaption).
		Category(mockCategory).
		Unit(mockUnit).
		TotalSteps(mockTotalSteps).
		Apply()

	// test via the getters that the builder values have been correctly applied.
	assert.Equal(t, task.GetName(), mockName)
	assert.Equal(t, task.GetCaption(), mockCaption)
	assert.Equal(t, task.GetCategory(), mockCategory)
	assert.Equal(t, task.GetUnit(), mockUnit)
	assert.Equal(t, task.GetTotalSteps(), mockTotalSteps)
}

type getterFunc[T comparable] func() T

type setterFunc[T comparable] func(T)

func testTaskGetterAndSetter[T comparable](t *testing.T, getter getterFunc[T], setter setterFunc[T], mockValue T) {
	t.Helper()

	var defaultValue T
	assert.NotEqual(t, defaultValue, mockValue, "setter is not adequately tested if mock value equals default value")
	assert.Equal(t, getter(), defaultValue)

	setter(mockValue)
	assert.Equal(t, getter(), mockValue)
}

// TestTaskName default, getter and setter work correctly.
func TestTaskName(t *testing.T) {
	task := createDefaultTask()
	testTaskGetterAndSetter(t, task.GetName, task.SetName, mockName)
}

// TestTaskCaption default, getter and setter work correctly.
func TestTaskCaption(t *testing.T) {
	task := createDefaultTask()
	testTaskGetterAndSetter(t, task.GetCaption, task.SetCaption, mockCaption)
}

// TestTaskCategory default, getter and setter work correctly.
func TestTaskCategory(t *testing.T) {
	task := createDefaultTask()
	testTaskGetterAndSetter(t, task.GetCategory, task.SetCategory, mockCategory)
}

// TestTaskUnit default unit and getter.
//
// Setting a custom unit is tested by [TestAddTask].
func TestTaskUnit(t *testing.T) {
	task := createDefaultTask()
	assert.Equal(t, &formatting.StepsUnit{}, task.GetUnit())
}

// TestTaskError default, getter and setter work correctly.
func TestTaskError(t *testing.T) {
	task := createDefaultTask()
	assert.False(t, task.IsError())

	// Setting a nil error is a no-op.
	task.Error(nil)
	assert.False(t, task.IsError())

	testTaskGetterAndSetter(t, task.GetError, task.Error, mockError)
	assert.True(t, task.IsError())

	// Once an error has been set, it cannot be cleared.
	task.Error(nil)
	assert.True(t, task.IsError())
}

func TestTask_GetStartedAt(t *testing.T) {
	beforeCreation := time.Now()

	// Ensure the task has a correctly assigned GetStartedAt time (by ensuring
	// it is a timestamp from after the creation time). The opposite is tested
	task := createDefaultTask()
	startedAt := task.GetStartedAt()

	assert.False(t, startedAt.IsZero())
	assert.GreaterOrEqual(t, startedAt, beforeCreation)

	// Assert that GetStartedAt hasn't changed.
	assert.Equal(t, task.GetStartedAt(), startedAt)
}

func TestTask_GetCompletedAt(t *testing.T) {
	task := createDefaultTask()
	assert.False(t, task.IsCompleted())
	assert.True(t, task.GetCompletedAt().IsZero())

	before := time.Now()
	task.CompleteStep()
	after := time.Now()
	assert.True(t, task.IsCompleted())
	assert.False(t, task.GetCompletedAt().IsZero())

	assert.GreaterOrEqual(t, task.GetCompletedAt(), before)
	assert.LessOrEqual(t, task.GetCompletedAt(), after)
}

func TestTask_GetElapsed(t *testing.T) {
	task := createDefaultTask()

	elapsed1 := task.GetElapsed()
	assert.Greater(t, elapsed1, time.Duration(0))

	elapsed2 := task.GetElapsed()
	assert.Greater(t, elapsed2, elapsed1)

	task.CompleteStep()
	assert.Equal(t, task.GetCompletedAt().Sub(task.GetStartedAt()), task.GetElapsed())
}

func TestTask_Indeterminate(t *testing.T) {
	task := createDefaultTask()
	assert.True(t, task.IsIndeterminate())
	assert.Equal(t, uint64(0), task.GetTotalSteps())

	assert.Equal(t, 0.0, task.GetProgress())
	assert.False(t, task.IsCompleted())
	assert.Equal(t, uint64(0), task.GetCompleteSteps())

	task.CompleteStep()
	assert.Equal(t, 1.0, task.GetProgress())
	assert.Equal(t, uint64(1), task.GetCompleteSteps())
	assert.True(t, task.IsCompleted())
}

func TestTask_Determinate(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().TotalSteps(2).Apply()
	assert.False(t, task.IsIndeterminate())
	assert.Equal(t, uint64(2), task.GetTotalSteps())

	assert.Equal(t, 0.0, task.GetProgress())
	assert.False(t, task.IsCompleted())
	assert.Equal(t, uint64(0), task.GetCompleteSteps())

	task.CompleteStep()
	assert.Equal(t, 0.5, task.GetProgress())
	assert.False(t, task.IsCompleted())
	assert.Equal(t, uint64(1), task.GetCompleteSteps())

	task.CompleteStep()
	assert.Equal(t, 1.0, task.GetProgress())
	assert.True(t, task.IsCompleted())
	assert.Equal(t, uint64(2), task.GetCompleteSteps())
}

func TestTask_TotalSteps(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().TotalSteps(2).Apply()

	assert.Equal(t, 0.0, task.GetProgress())
	assert.False(t, task.IsCompleted())
	assert.Equal(t, uint64(0), task.GetCompleteSteps())

	task.CompleteSteps(1)
	assert.Equal(t, 0.5, task.GetProgress())
	assert.False(t, task.IsCompleted())
	assert.Equal(t, uint64(1), task.GetCompleteSteps())

	// Assert that changing the number of total steps changes the progress
	// automatically.
	task.TotalSteps(4)
	assert.Equal(t, 0.25, task.GetProgress())
	assert.False(t, task.IsCompleted())
	assert.Equal(t, uint64(1), task.GetCompleteSteps())

	// Complete the task
	task.CompleteSteps(3)
	assert.Equal(t, 1.0, task.GetProgress())
	assert.True(t, task.IsCompleted())
	assert.Equal(t, uint64(4), task.GetCompleteSteps())

	// Assert indirectly that this is a no-op.
	task.TotalSteps(100)
	assert.Equal(t, 1.0, task.GetProgress())
	assert.True(t, task.IsCompleted())
	assert.Equal(t, uint64(4), task.GetCompleteSteps())
}

func TestTask_CompleteSteps_zero_is_no_op(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().TotalSteps(2).Apply()
	assert.Equal(t, uint64(0), task.GetCompleteSteps())

	// Assert that this is a no-op.
	task.CompleteSteps(0)
	assert.Equal(t, uint64(0), task.GetCompleteSteps())
}

func TestTask_CompleteSteps_completed_is_no_op(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().TotalSteps(4).Apply()
	assert.Equal(t, uint64(0), task.GetCompleteSteps())
	assert.False(t, task.IsCompleted())

	task.CompleteSteps(4)
	assert.Equal(t, uint64(4), task.GetCompleteSteps())
	assert.True(t, task.IsCompleted())

	// Assert that this is a no-op.
	task.CompleteSteps(100)
	assert.Equal(t, uint64(4), task.GetCompleteSteps())
	assert.True(t, task.IsCompleted())
}

func TestTask_CompletedSteps_clamp(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().TotalSteps(2).Apply()
	assert.Equal(t, uint64(0), task.GetCompleteSteps())

	// The three completed steps should be clamped to the two total steps.
	task.CompleteSteps(3)
	assert.Equal(t, uint64(2), task.GetCompleteSteps())
}

func TestTask_SetCompletedSteps_zero_is_no_op(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().TotalSteps(2).Apply()
	assert.Equal(t, uint64(0), task.GetCompleteSteps())

	// Assert that this is a no-op.
	task.SetCompletedSteps(0)
	assert.Equal(t, uint64(0), task.GetCompleteSteps())
}

func TestTask_SetCompletedSteps_less_is_no_op(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().TotalSteps(4).Apply()
	assert.Equal(t, uint64(0), task.GetCompleteSteps())

	task.SetCompletedSteps(2)
	assert.Equal(t, 0.5, task.GetProgress())
	assert.False(t, task.IsCompleted())
	assert.Equal(t, uint64(2), task.GetCompleteSteps())

	task.SetCompletedSteps(1)
	assert.Equal(t, 0.5, task.GetProgress())
	assert.False(t, task.IsCompleted())
	assert.Equal(t, uint64(2), task.GetCompleteSteps())
}

func TestTask_SetCompletedSteps_completed_is_no_op(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().TotalSteps(4).Apply()
	assert.Equal(t, uint64(0), task.GetCompleteSteps())
	assert.False(t, task.IsCompleted())

	task.SetCompletedSteps(4)
	assert.Equal(t, uint64(4), task.GetCompleteSteps())
	assert.True(t, task.IsCompleted())

	// Assert that this is a no-op.
	task.SetCompletedSteps(100)
	assert.Equal(t, uint64(4), task.GetCompleteSteps())
	assert.True(t, task.IsCompleted())
}

func TestTask_SetCompletedSteps_clamp(t *testing.T) {
	m := mon.New("test")
	task := m.AddTask().TotalSteps(2).Apply()
	assert.Equal(t, uint64(0), task.GetCompleteSteps())

	// The three completed steps should be clamped to the two total steps.
	task.SetCompletedSteps(3)
	assert.Equal(t, uint64(2), task.GetCompleteSteps())
}
