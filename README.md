# mon(itor)

[![Go Reference](https://pkg.go.dev/badge/github.com/apollosoftwarexyz/mon.svg)](https://pkg.go.dev/github.com/apollosoftwarexyz/mon)

A simple and clean abstraction for reporting progress to a user in a CLI.

<img width="80%" src="./demo/demo.gif" alt="go run ./demo/main.go">

## Features

- Built-in CLI animations and terminal UI.
- Runs asynchronously in the background using goroutines.
- Fluent builder API for adding tasks to the monitor.
- Supports determinate (fixed number of steps) and indeterminate tasks.
- Automatic estimated completion time calculations for determinate tasks.
- Extensible human-readable [formatting](https://pkg.go.dev/github.com/apollosoftwarexyz/mon/formatting) for values.

## Usage

Import the [`mon`](https://pkg.go.dev/github.com/apollosoftwarexyz/mon) library
and create a new monitor.

```go
package main

import (
	"context"

	"github.com/apollosoftwarexyz/mon"
)
```

Then, create and show a new monitor. As soon as `Show` is called, the monitor is
displayed in the terminal. It can be removed by calling the returned `cancel()`
function which should always be called with a `defer`.

```go
// Create a new monitor.
m := mon.New("Please wait")

// Show the monitor. The returned cancel method allows the UI to be cleaned
// up programmatically or automatically (with defer) when the work is done.
ctx, cancel := m.Show(context.WithCancelCause(context.Background()))
defer cancel(nil)
```

Finally, track some tasks! You can get a new task builder by calling `AddTask`
on the monitor.

- Use `Apply` on the builder when you want to show the task on the monitor
  (normally this is immediately).
- For determinate tasks (where the number of total steps is known), set
  `TotalSteps` on the builder as shown below and then track work with
  `CompleteStep` or `CompleteSteps`.
- As steps are completed, statistics such as the estimated remaining time and
  number of completed steps per second are automatically computed.

```go
// Do the work
task := m.AddTask().Name("my task").TotalSteps(2048).Apply()

// Increment the number of completed steps.
task.CompleteStep()
task.CompleteStep()

// Bulk increment the number of completed steps.
task.CompleteSteps(1024)
```

Indeterminate tasks are also supported, simply don't set `TotalSteps` and call
`CompleteStep` when done:

```go
// Also supports indeterminate tasks!
indeterminateTask := m.AddTask().Name("indeterminate task").Apply()

// ...just use CompleteStep when it's done!
indeterminateTask.CompleteStep()
```