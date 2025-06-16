# Extending Gouri

This guide explains how to add new commands or modify existing ones.

1. Open `main.go` and create a helper function implementing the feature.
2. Add a `case` statement inside `main()` for the new command that calls your helper.
3. Update the `usage()` function so the command appears in the short help output.
4. Append a description to the `manualText` constant so it shows up in `gouri manual`.
5. If the command requires persistent settings, modify the `Config` struct in `main.go`.
6. Run `go build` and `go vet ./...` to ensure the code compiles and passes basic linting.
7. Update `README.md` and `docs/manual.md` with usage information.

By keeping the manual, usage text and documentation in sync, you ensure users
understand the new functionality.
