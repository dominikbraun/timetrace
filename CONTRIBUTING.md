# Contributing

## Feature requests

Desire a new timetrace feature? Just propose your idea by
[creating an issue](https://github.com/dominikbraun/timetrace/issues/new).

Good feature proposals ...
* explain the problem that the feature solves.
* explain why it would be a desirable feature for the majority of users.

## Reporting issues

If you encounter an unexpected behavior or a bug, feel free to
[file an issue](https://github.com/dominikbraun/timetrace/issues/new). When you
do so, please make sure to ...
* include version information from the output of `timetrace version`.
* provide steps to reproduce the behavior.

## Code contributions

### Setting up local development

Developing timetrace only requires [Go 1.16](https://golang.org/dl/).

1. Fork the repository.
2. Clone your forked repository.
3. Run `go run . version` to verify that everything works.

### Coding conventions

* All code has to follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.
* All code has to be formatted with `gofmt -s`.
* Exported types and methods should be documented briefly. Explain what they're doing, not what they are.

### Contributing changes

1. Make your changes. If needed, [write tests](core/timetrace_test.go).
2. Run `go run . <command>` for testing your changes.
3. Run `go test ./...` to verify that all tests pass.
4. Commit your changes and open a PR.
