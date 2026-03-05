package mon

import (
	"context"
	"time"

	"github.com/apollosoftwarexyz/mon/animations"
	tea "github.com/charmbracelet/bubbletea"
)

// M is a CLI monitor for various [Task] statuses.
type M interface {
	// AddTask creates a [TaskBuilder] that can be used to define and add a new
	// task to the monitor.
	AddTask() TaskBuilder

	// IsCancellationBlocked returns true if cancellation has been blocked with
	// [M.BlockCancellation] (or if it has been unblocked with
	// [M.AllowCancellation]).
	IsCancellationBlocked() bool

	// AllowCancellation permits the user to press Ctrl+C to cancel the context
	// used by the monitor.
	//
	// This is the default behavior and need only be called to restore
	// cancellation behavior if it has been blocked with [M.BlockCancellation].
	//
	// The same monitor instance is returned to allow for a fluent API.
	AllowCancellation() M

	// BlockCancellation ignores user actions such as Ctrl+C that would
	// otherwise cancel the context used by the monitor.
	//
	// This function is unable to block all forms of cancellation (for example,
	// SIGKILL) - the intention is just to help discourage the user from
	// cancelling a task at a bad moment.
	//
	// The same monitor instance is returned to allow for a fluent API.
	BlockCancellation() M

	// GetCaption of the monitor.
	GetCaption() string

	// SetCaption of the monitor.
	SetCaption(caption string)

	// Show the monitor in the CLI.
	//
	// The [CancelFunc] should be deferred immediately after Show is called:
	//
	//	ctx, cancel := m.Show(context.WithCancelCause(context.Background()))
	//	defer cancel()
	Show(ctx context.Context, cancel context.CancelCauseFunc) (context.Context, context.CancelCauseFunc)
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
		exited:       make(chan error),
	}
}

func (m *model) AddTask() TaskBuilder {
	return &taskBuilder{m: m}
}

func (m *model) IsCancellationBlocked() bool {
	return m.blockCancellation
}

func (m *model) AllowCancellation() M {
	m.blockCancellation = false
	return m
}

func (m *model) BlockCancellation() M {
	m.blockCancellation = true
	return m
}

func (m *model) GetCaption() string {
	return m.caption
}

func (m *model) SetCaption(caption string) {
	m.caption = caption
	m.notify()
}

func (m *model) Show(ctx context.Context, cancel context.CancelCauseFunc) (context.Context, context.CancelCauseFunc) {
	m.prog = tea.NewProgram(m, tea.WithContext(ctx))

	go func() {
		_, err := m.prog.Run()
		cancel(err)
		m.exited <- err
		close(m.exited)
	}()

	return ctx, func(cause error) {
		m.notifyDone()

		// Wait for the bubbletea application to quit.
		select {
		case <-m.exited:
		}
	}
}
