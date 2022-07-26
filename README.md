# nakama

A Go web and realtime for Nakama game server. Works with wasm builds.

[![Tests](https://github.com/ascii8/nakama-go/workflows/Test/badge.svg)](https://github.com/ascii8/nakama-go/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/ascii8/nakama-go)](https://goreportcard.com/report/github.com/ascii8/nakama-go)
[![Reference](https://godoc.org/github.com/ascii8/nakama-go?status.svg)](https://pkg.go.dev/github.com/ascii8/nakama-go)
[![Releases](https://img.shields.io/github/v/release/ascii8/nakama-go?display_name=tag&sort=semver)](https://github.com/ascii8/nakama-go/releases)

## using

```sh
go get github.com/ascii8/nakama-go
```

## quickstart

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

## examples

See the [Go package documentation](https://pkg.go.dev/github.com/ascii8/nakama-go)
for more examples.


## related

* [`github.com/ascii8/nktest`](https://github.com/ascii8/nktest) - a Nakama test runner.
* [`github.com/ascii8/xoxo-go`](https://github.com/ascii8/xoxo-go) - a Go version of Nakama's phaserjs xoxo example. Uses Ebitengine
