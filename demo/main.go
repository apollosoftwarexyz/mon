package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apollosoftwarexyz/mon"
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
	countToN(&wg, m, 20000)
	countToN(&wg, m, 5000)
	countToN(&wg, m, 2000)
	countToN(&wg, m, 2000)
	countToN(&wg, m, 500)
	countToN(&wg, m, 1000)
	countToN(&wg, m, 300)
	countToN(&wg, m, 200)
	wg.Wait()

}

func indeterminateTask(wg *sync.WaitGroup, m mon.M, duration time.Duration) {

	task := m.AddTask().Name("mysterious task").Apply()

	wg.Go(func() {
		time.Sleep(duration)
		task.CompletedStep()
	})

}

func countToN(wg *sync.WaitGroup, m mon.M, n uint64) {

	task := m.AddTask().Name(fmt.Sprintf("counting to %d", n)).TotalSteps(n).Apply()

	wg.Go(func() {
		for range n {
			task.CompletedStep()
			time.Sleep(time.Millisecond)
		}
	})

}
