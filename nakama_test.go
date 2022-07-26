package nakama

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/ascii8/nktest"
)

// globalCtx is the global context.
var globalCtx context.Context

// nkTest is the nakama test runner.
var nkTest *nktest.Runner

func TestHealthcheck(t *testing.T) {
	cl := New(WithURL(nkTest.HttpLocal()))
	if err := cl.Healthcheck(globalCtx); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRpc(t *testing.T) {
	cl := New(WithURL(nkTest.HttpLocal()))
	var res rewards
	amount := int64(1000)
	if err := Rpc("dailyRewards").
		WithHttpKey(nkTest.Name()).
		WithPayload(rewards{
			Rewards: amount,
		}).Do(globalCtx, cl, &res); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("rewards: %d", res.Rewards)
	if res.Rewards != 2*amount {
		t.Errorf("expected %d, got: %d", 2*amount, res.Rewards)
	}
}

type rewards struct {
	Rewards int64 `json:"rewards,omitempty"`
}

// TestMain handles setting up and tearing down the postgres and nakama docker
// images.
func TestMain(m *testing.M) {
	var cancel func()
	globalCtx, cancel = context.WithCancel(context.Background())
	go func() {
		// catch signals, canceling context to cause cleanup
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-globalCtx.Done():
		case sig := <-ch:
			fmt.Fprintf(os.Stdout, "SIGNAL: %s\n", sig)
			cancel()
		}
	}()
	code := 0
	pull := os.Getenv("PULL")
	nkTest = nktest.New(
		nktest.WithAlwaysPull(pull != "" && pull != "false" && pull != "0"),
		nktest.WithBuildConfig("./apitest", nktest.WithDefaultGoEnv(), nktest.WithDefaultGoVolumes()),
	)
	if err := nkTest.Run(globalCtx); err == nil {
		code = m.Run()
	} else {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		code = 1
	}
	cancel()
	<-time.After(2200 * time.Millisecond)
	os.Exit(code)
}
