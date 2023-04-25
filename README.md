# About

A Go web and realtime client package for the Nakama game server. Works with
WASM builds.

[![Tests](https://github.com/ascii8/nakama-go/workflows/Test/badge.svg)](https://github.com/ascii8/nakama-go/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/ascii8/nakama-go)](https://goreportcard.com/report/github.com/ascii8/nakama-go)
[![Reference](https://pkg.go.dev/badge/github.com/ascii8/nakama-go.svg)](https://pkg.go.dev/github.com/ascii8/nakama-go)
[![Releases](https://img.shields.io/github/v/release/ascii8/nakama-go?display_name=tag&sort=semver)](https://github.com/ascii8/nakama-go/releases)

## Using

```sh
go get github.com/ascii8/nakama-go
```

## Quickstart

```go
package nakama_test

import (
	"context"
	"fmt"
	"log"

	"github.com/ascii8/nakama-go"
)

func Example() {
	const id = "6d0c9e83-8385-48a8-8601-060b8f6a3bf6"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// create client
	cl := nakama.New(nakama.WithServerKey("apitest_server"))
	// authenticate
	if err := cl.AuthenticateDevice(ctx, id, true, ""); err != nil {
		log.Fatal(err)
	}
	// retrieve account
	account, err := cl.Account(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// list devices on the account
	for _, d := range account.Devices {
		fmt.Println("id:", d.Id)
	}
	// Output:
	// id: 6d0c9e83-8385-48a8-8601-060b8f6a3bf6
}

```

## Examples

See [`github.com/ascii8/xoxo-go`](https://github.com/ascii8/xoxo-go) for a
demonstration of end-to-end unit tests using `nktest`, and "real world"
examples of pure Go clients (e.g., Ebitengine) built using this client
package.

See the [Go package documentation](https://pkg.go.dev/github.com/ascii8/nakama-go)
for other examples.

## Notes

Run browser tests:

```sh
# setup wasmbrowsertest
$ go install github.com/agnivade/wasmbrowsertest@latest
$ cd $GOPATH/bin && ln -s wasmbrowsertest go_js_wasm_exec

# run the wasmtests
$ cd /path/to/nakama-go/wasmtest
$ GOOS=js GOARCH=wasm go test -v
```

## Related Links

* [`github.com/ascii8/nktest`](https://github.com/ascii8/nktest) - a Nakama test runner
* [`github.com/ascii8/xoxo-go`](https://github.com/ascii8/xoxo-go) - a pure Go version of Nakama's XOXO example, demonstrating end-to-end unit tests, and providing multiple example clients using this package. Has example Ebitengine client
