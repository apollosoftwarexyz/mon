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
	cancel := m.Show(context.Background())
	defer cancel()

	// Do the work...
	var wg sync.WaitGroup

	indeterminateTask(&wg, m, 10*time.Second)
	indeterminateTask(&wg, m, 5*time.Second)
	indeterminateTask(&wg, m, 200*time.Millisecond)
	errorTask(&wg, m, 5*time.Second)

	fakeCopyBytes(&wg, m, 20000)
	fakeCopyBytes(&wg, m, 5000)
	fakeCopyBytes(&wg, m, 2000)
	fakeCopyBytes(&wg, m, 2000)
	fakeCopyBytes(&wg, m, 500)
	fakeCopyBytes(&wg, m, 1000)
	fakeCopyBytes(&wg, m, 300)
	fakeCopyBytes(&wg, m, 200)

	time.Sleep(3 * time.Second)
	fakeCopyBytes(&wg, m, 2000)
	fakeCopyBytes(&wg, m, 15000)

	time.Sleep(5 * time.Second)
	fakeCopyBytes(&wg, m, 2000)

	wg.Wait()

}

func indeterminateTask(wg *sync.WaitGroup, m mon.M, duration time.Duration) {

	task := m.AddTask().Name("mysterious task").Apply()

	wg.Go(func() {
		time.Sleep(duration)
		task.CompleteStep()
	})

}

func fakeCopyBytes(wg *sync.WaitGroup, m mon.M, n uint64) {

	task := m.AddTask().
		Name(fmt.Sprintf("copying %d bytes", n)).
		Unit(&formatting.BytesUnit{}).
		TotalSteps(n).
		Apply()

	wg.Go(func() {
		for range n {
			task.CompleteStep()
			time.Sleep(time.Millisecond)
		}
	})

}

func errorTask(wg *sync.WaitGroup, m mon.M, duration time.Duration) {

	task := m.AddTask().Name("risky task").Apply()

	wg.Go(func() {
		time.Sleep(duration)
		task.Error(errors.New("this is a simulated error"))
	})

}
