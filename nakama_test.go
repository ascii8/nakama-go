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
	ctx, cancel := context.WithCancel(globalCtx)
	defer cancel()
	cl := newClient(ctx, t, false)
	if err := cl.Healthcheck(ctx); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRpc(t *testing.T) {
	ctx, cancel := context.WithCancel(globalCtx)
	defer cancel()
	cl := newClient(ctx, t, false)
	var res rewards
	amount := int64(1000)
	if err := Rpc("dailyRewards").
		WithHttpKey(nkTest.Name()).
		WithPayload(rewards{
			Rewards: amount,
		}).Do(ctx, cl, &res); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("rewards: %d", res.Rewards)
	if res.Rewards != 2*amount {
		t.Errorf("expected %d, got: %d", 2*amount, res.Rewards)
	}
}

func TestWebsocket(t *testing.T) {
	ctx, cancel := context.WithCancel(globalCtx)
	defer cancel()
	conn, err := newClient(ctx, t, true).Dial(ctx, WithDialCreateStatus(true))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	defer conn.Close()
}

type rewards struct {
	Rewards int64 `json:"rewards,omitempty"`
}

func newClient(ctx context.Context, t *testing.T, addProxyLogger bool, opts ...Option) *Client {
	local := nkTest.HttpLocal()
	t.Logf("real: %s", local)
	logger := nktest.NewLogger(t.Logf)
	var proxyOpts []nktest.ProxyOption
	if addProxyLogger {
		proxyOpts = append(proxyOpts, nktest.WithLogger(logger))
	}
	urlstr, err := nktest.NewProxy(proxyOpts...).Run(ctx, local)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("url: %s", urlstr)
	return New(append([]Option{
		WithURL(urlstr),
		WithTransport(logger.Transport(nil)),
	}, opts...)...)
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
		nktest.WithDir("./apitest"),
		nktest.WithAlwaysPull(pull != "" && pull != "false" && pull != "0"),
		nktest.WithBuildConfig("./nkapitest", nktest.WithDefaultGoEnv(), nktest.WithDefaultGoVolumes()),
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
