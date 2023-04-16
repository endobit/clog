[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/endobit/clog/test.yaml)](https://github.com/endobit/clog/actions?query=workflow%3Atest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/endobit/clog)](https://img.shields.io/github/go-mod/go-version/endobit/clog)
[![Go Report Card](https://goreportcard.com/badge/github.com/endobit/clog)](https://goreportcard.com/report/github.com/endobit/clog)
[![Codecov](https://codecov.io/gh/endobit/oui/branch/main/graph/badge.svg)](https://codecov.io/gh/endobit/clog)
[![Go Reference](https://pkg.go.dev/badge/github.com/endobit/clog.svg)](https://pkg.go.dev/github.com/endobit/clog)

# Clog

Color logging with
[golang.org/x/exp/slog](https://pkg.go.dev/golang.org/x/exp/slog). Clog mimics
[`zerolog.ConsoleWriter`](https://github.com/rs/zerolog#readme) style but due to
the `slog.Handler` implementation field order is preserved, whereas the
`zerolog.ConsoleWriter` parses the json and sorts the field.

![Logging Sample](sample.png)

## Stability

Clog will track the `slog` package including any breaking changes.








