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
	"github.com/google/uuid"
	nkapi "github.com/heroiclabs/nakama-common/api"
	"golang.org/x/exp/slices"
)

// globalCtx is the global context.
var globalCtx context.Context

// nkTest is the nakama test runner.
var nkTest *nktest.Runner

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

func TestAuthenticateDevice(t *testing.T) {
	ctx, cancel := context.WithCancel(globalCtx)
	defer cancel()
	cl := newClient(ctx, t, false, WithServerKey(nkTest.ServerKey()))
	createAccount(ctx, t, cl)
}

func TestWebsocket(t *testing.T) {
	ctx, cancel := context.WithCancel(globalCtx)
	defer cancel()
	cl := newClient(ctx, t, false, WithServerKey(nkTest.ServerKey()))
	createAccount(ctx, t, cl)
	conn, err := cl.NewConn(ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	defer conn.Close()
	select {
	case <-ctx.Done():
	case <-time.After(3 * time.Minute):
	}
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
		WithServerKey(nkTest.ServerKey()),
	}, opts...)...)
}

func createAccount(ctx context.Context, t *testing.T, cl *Client) {
	deviceId := uuid.New().String()
	t.Logf("registering: %s", deviceId)
	if err := cl.AuthenticateDevice(ctx, true, deviceId, ""); err != nil {
		t.Fatalf("expected no error: got: %v", err)
	}
	expiry := cl.SessionExpiry()
	t.Logf("expiry: %s", cl.SessionExpiry())
	if expiry.IsZero() || expiry.Before(time.Now()) {
		t.Fatalf("expected non-zero expiry in the future, got: %s", expiry)
	}
	res, err := cl.Account(ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("account: %+v", res)
	if len(res.Devices) == 0 {
		t.Fatalf("expected there to be at least one device")
	}
	i := slices.IndexFunc(res.Devices, func(d *nkapi.AccountDevice) bool {
		return d.Id == deviceId
	})
	if i == -1 {
		t.Fatalf("expected accountRes.Devices to contain %s", deviceId)
	}
}
