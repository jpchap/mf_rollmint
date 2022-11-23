# rollmint

ABCI-client implementation for Optimistic Rollups.

This fork specifically includes ABCI++ Vote Extensions. These will be used to decentralize the oracle set and involve them in block production.

Design document: <https://docs.google.com/document/d/12gZow_JTJjRrmaD2mNTmYniLhyxVLSyDd7Fbxo5UnA8/edit?usp=sharing>

[![build-and-test](https://github.com/celestiaorg/rollmint/actions/workflows/test.yml/badge.svg)](https://github.com/celestiaorg/rollmint/actions/workflows/test.yml)
[![golangci-lint](https://github.com/celestiaorg/rollmint/actions/workflows/lint.yml/badge.svg)](https://github.com/celestiaorg/rollmint/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/celestiaorg/rollmint)](https://goreportcard.com/report/github.com/celestiaorg/rollmint)
[![codecov](https://codecov.io/gh/celestiaorg/rollmint/branch/main/graph/badge.svg?token=CWGA4RLDS9)](https://codecov.io/gh/celestiaorg/rollmint)
[![GoDoc](https://godoc.org/github.com/celestiaorg/rollmint?status.svg)](https://godoc.org/github.com/celestiaorg/rollmint)
[![Twitter Follow](https://img.shields.io/twitter/follow/CelestiaOrg?style=social)](https://twitter.com/CelestiaOrg)

## Building From Source

Requires Go version >= 1.17.

To build:

```sh
git clone https://github.com/celestiaorg/rollmint.git
cd rollmint
go build -v ./...
```

To test:

```sh
go test ./...
```

To regenerate protobuf types:

```sh
./proto/gen.sh
```

## Contributing

We welcome your contributions! Everyone is welcome to contribute, whether it's in the form of code,
documentation, bug reports, feature requests, or anything else.

If you're looking for issues to work on, try looking at the [good first issue list](https://github.com/celestiaorg/rollmint/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22). Issues with this tag are suitable for a new external contributor and is a great way to find something you can help with!

See [the contributing guide](./CONTRIBUTING.md) for more details.

Please join our [Community Discord](https://discord.com/invite/YsnTPcSfWQ) to ask questions, discuss your ideas, and connect with other contributors.

## Code of Conduct

See our Code of Conduct [here](https://docs.celestia.org/community/coc).
