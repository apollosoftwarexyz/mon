package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/apollosoftwarexyz/mon"
	"github.com/apollosoftwarexyz/mon/formatting"
)

func main() {

	// Create a new monitor.
	m := mon.New("Please wait")

	// Show the monitor. The returned cancel method allows the UI to be cleaned
	// up automatically when the work is done.
	ctx, cancel := m.Show(context.WithCancelCause(context.Background()))
	defer cancel(nil)

	// Do some "work"...
	var wg sync.WaitGroup

	indeterminateTask(&wg, ctx, m, 10*time.Second)
	indeterminateTask(&wg, ctx, m, 5*time.Second)
	indeterminateTask(&wg, ctx, m, 200*time.Millisecond)

	errorTask(&wg, ctx, m, 5*time.Second)

	fakeCopyBytes(&wg, ctx, m, 5000)
	fakeCopyBytes(&wg, ctx, m, 2000)
	fakeCopyBytes(&wg, ctx, m, 2000)
	fakeCopyBytes(&wg, ctx, m, 500)
	fakeCopyBytes(&wg, ctx, m, 1000)
	fakeCopyBytes(&wg, ctx, m, 300)
	fakeCopyBytes(&wg, ctx, m, 200)

	interruptableSleep(ctx, 3*time.Second)
	fakeCopyBytes(&wg, ctx, m, 2000)

	interruptableSleep(ctx, 5*time.Second)
	fakeCopyBytes(&wg, ctx, m, 2000)

	interruptableSleep(ctx, 4*time.Second)
	wg.Wait()

}

func indeterminateTask(wg *sync.WaitGroup, ctx context.Context, m mon.M, duration time.Duration) {

	task := m.AddTask().Name("mysterious task").Apply()

	wg.Go(func() {
		interruptableSleep(ctx, duration)
		task.CompleteStep()
	})

}

func fakeCopyBytes(wg *sync.WaitGroup, ctx context.Context, m mon.M, n uint64) {

	task := m.AddTask().
		Name(fmt.Sprintf("copying %d bytes", n)).
		Unit(&formatting.BytesUnit{}).
		TotalSteps(n).
		Apply()

	wg.Go(func() {
		for range n {
			task.CompleteStep()
			interruptableSleep(ctx, time.Millisecond)
		}
	})

}

func errorTask(wg *sync.WaitGroup, ctx context.Context, m mon.M, duration time.Duration) {

	task := m.AddTask().Name("risky task").Apply()

	wg.Go(func() {
		interruptableSleep(ctx, duration)
		task.Error(errors.New("this is a simulated error"))
	})

}

func interruptableSleep(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
}
