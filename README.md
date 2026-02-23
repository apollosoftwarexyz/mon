# mon(itor)

[![Go Reference](https://pkg.go.dev/badge/github.com/apollosoftwarexyz/mon.svg)](https://pkg.go.dev/github.com/apollosoftwarexyz/mon)

A progress monitor for Go applications.

This package aims to provide a simple and clean abstraction for reporting
progress to a user in a CLI.

<div style="text-align: center;">
    <img width="80%" src="./demo/demo.gif" alt="go run ./demo/main.go">
</div>

## Features

- Built-in CLI animations and terminal UI.
- Runs asynchronously in the background using goroutines.
- Fluent builder API for adding tasks to the monitor.
- Supports determinate (fixed number of steps) and indeterminate tasks.
- Automatic estimated completion time calculations for determinate tasks.
- Extensible human-readable [formatting](https://pkg.go.dev/github.com/apollosoftwarexyz/mon/formatting) for values.
