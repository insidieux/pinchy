# Pinchy

Service discovering and registry bridge.

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/insidieux/pinchy/CI?style=flat-square)](https://github.com/insidieux/pinchy/actions?query=workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/insidieux/pinchy)](https://goreportcard.com/report/github.com/insidieux/pinchy)
[![codecov](https://codecov.io/gh/insidieux/pinchy/branch/master/graph/badge.svg?token=BI6HEMPLB1)](https://codecov.io/gh/insidieux/pinchy/branch/master)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/insidieux/pinchy)

Pinchy is a simple binary, which allows to automatically fetch services info from `Source` and register/remove them to/from `Registry`.

Supported pluggable service sources:
- [YAML File](https://ru.wikipedia.org/wiki/YAML)

Supported pluggable service registries:
- [Consul](http://www.consul.io/)

## Installing

Install Pinchy by running:

```shell
go get github.com/insidieux/pinchy/cmd/pinchy
```

Ensure that `$GOPATH/bin` is added to your `$PATH`.

## Documentation

- [User guide][]
- [Contributing guide][]

[User guide]: ./docs/user-guide.md
[Contributing guide]: ./docs/contributing.md

## License

[Apache][]

[Apache]: ./LICENSE
