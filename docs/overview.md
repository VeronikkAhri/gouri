# Overview

Gouri is a command line assistant written in Go. It aims to simplify everyday
shell tasks with short commands. The tool works on Linux and Windows
PowerShell, abstracting away platform differences when possible.

The application is a single binary. Build it using `go build` and then run the
`gouri` executable to execute commands. Most operations simply wrap standard
system utilities, providing a unified interface that saves typing and removes
complex flags.

Gouri can also store user preferences in `~/.gouri.json`, letting you define a
default editor or custom key/value pairs for later retrieval.

Run `gouri manual` to display a full list of commands with their descriptions.
