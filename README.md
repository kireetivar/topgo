# Topgo

Topgo is a simple and efficient Linux system monitor written in Go. It relies on the Linux `/proc` filesystem to gather system and process-level metrics (CPU, memory, etc.) and presents them in a rich Terminal User Interface (TUI) using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

## Requirements

- **Go**: Version 1.24.2 or later.
- **Operating System**: Linux (or a compatible environment like WSL) is required, as the application directly reads from the `/proc` filesystem.

## Building and Running

A `Makefile` is included to simplify the build and execution process.

To format, vet, and build the executable:
```bash
make build
```

To format and run the application directly:
```bash
make run
```

Alternatively, using standard Go commands:
```bash
go build -o topgo main.go
./topgo
```

## Features
- Real-time system monitoring.
- Per-process CPU and memory tracking.
