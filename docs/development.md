# Development

Gouri is built with Go 1.22 or later. The repository includes a `go.mod` file
so you can simply run `go build` to compile. Linting can be performed with
`go vet ./...`.

To add features, edit `main.go`. The file contains helper functions for each
command and a large switch statement that dispatches based on the first
argument. When adding new commands, update the `usage()` function and the
`manualText` constant so the documentation stays consistent.

Pull requests should include tests where appropriate and keep the
documentation in `README.md` and the `docs/` folder up to date.
