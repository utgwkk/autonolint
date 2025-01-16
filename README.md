# autonolint

Insert `//nolint` comment automatically for golangci-lint

## Installation

```console
$ go install github.com/utgwkk/autonolint/cmd/autonolint
```

## Usage

1. Run golangci-lint with `--out-format=json`.
2. Extract `Issues` field of JSON.
3. Pass to `autonolint` with standard input.

```console
$ golangci-lint run --out-format=json | jq .Issues | autonolint
```

You can specify content of a comment for `//nolint` directive via `-comment` option.

```console
$ golangci-lint run --out-format=json | jq .Issues | autonolint -comment "TODO: refactor"
```

## DEMO

[![Image from Gyazo](https://i.gyazo.com/82570c43fcaf57c6dcd9eff93fcb8b8e.gif)](https://gyazo.com/82570c43fcaf57c6dcd9eff93fcb8b8e)
