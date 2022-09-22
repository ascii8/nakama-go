package nakama_test

import (
	"context"
	"fmt"
	"log"

	"github.com/ascii8/nakama-go"
)

// , nakama.WithTransport(nktest.NewLogger(log.Printf).Transport(nil))

func Example() {
	const id = "6d0c9e83-8385-48a8-8601-060b8f6a3bf6"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// create client
	cl := nakama.New(nakama.WithServerKey("testdata_server"))
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

func ExampleRpc() {
	const amount = 1000
	type rewards struct {
		Rewards int64 `json:"rewards,omitempty"`
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// create client
	cl := nakama.New(nakama.WithServerKey("testdata_server"))
	// create request and response
	res := new(rewards)
	req := nakama.Rpc("dailyRewards", rewards{Rewards: amount}, res)
	// execute rpc with http key
	if err := req.WithHttpKey("testdata").Do(ctx, cl); err != nil {
		log.Fatal(err)
	}
	fmt.Println("rewards:", res.Rewards)
	// Output:
	// rewards: 2000
}
