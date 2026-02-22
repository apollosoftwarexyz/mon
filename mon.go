package mon

import (
	"context"
	"time"

	"github.com/apollosoftwarexyz/mon/animations"
	tea "github.com/charmbracelet/bubbletea"
)

// CancelFunc automatically cancels the [context.Context] that was passed to
// [M.Show] and waits for the monitor to be cleaned up.
type CancelFunc func()

// M is a CLI monitor for various [Task] statuses.
type M interface {
	// AddTask creates a [TaskBuilder] that can be used to define and add a new
	// task to the monitor.
	AddTask() TaskBuilder

	// Show the monitor in the CLI.
	//
	// The [CancelFunc] should be deferred immediately after Show is called:
	//
	//	cancel := m.Show(context.Background())
	//	defer cancel()
	Show(ctx context.Context) CancelFunc
}

// New monitor.
//
// Once configured, the monitor can be displayed with [M.Show].
func New(caption string) M {
	return &model{
		spinnerAnim:  animations.Default(),
		ellipsisAnim: animations.Ellipsis(),
		start:        time.Now(),
		caption:      caption,
	}
}

func (m *model) AddTask() TaskBuilder {
	return &taskBuilder{m: m}
}

func (m *model) Show(ctx context.Context) CancelFunc {
	m.prog = tea.NewProgram(m, tea.WithContext(ctx))
	exited := make(chan error)

	go func() {
		_, err := m.prog.Run()
		exited <- err
	}()

	return func() {
		m.prog.Send(doneMsg{})

		// Wait for the bubbletea application to quit.
		select {
		case <-exited:
		}
	}
}
