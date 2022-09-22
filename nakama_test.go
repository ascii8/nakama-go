package nakama

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ascii8/nktest"
	"github.com/google/uuid"
	nkapi "github.com/heroiclabs/nakama-common/api"
	"golang.org/x/exp/slices"
)

// TestMain handles setting up and tearing down the postgres and nakama
// containers.
func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx = nktest.WithAlwaysPullFromEnv(ctx, "PULL")
	ctx = nktest.WithHostPortMap(ctx)
	nktest.Main(ctx, m,
		nktest.WithDir("./testdata"),
		nktest.WithBuildConfig("./nkapitest", nktest.WithDefaultGoEnv(), nktest.WithDefaultGoVolumes()),
	)
}

func TestHealthcheck(t *testing.T) {
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	defer cancel()
	cl := newClient(ctx, t, nk)
	if err := cl.Healthcheck(ctx); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRpc(t *testing.T) {
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	defer cancel()
	const amount int64 = 1000
	cl := newClient(ctx, t, nk)
	var res rewards
	req := Rpc(
		"dailyRewards",
		rewards{
			Rewards: amount,
		},
		&res,
	).
		WithHttpKey(nk.Name())
	if err := req.Do(ctx, cl); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("rewards: %d", res.Rewards)
	if res.Rewards != 2*amount {
		t.Errorf("expected %d, got: %d", 2*amount, res.Rewards)
	}
}

func TestRpcProtoEncodeDecode(t *testing.T) {
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	defer cancel()
	const name string = "bob"
	const amount int64 = 1000
	cl := newClient(ctx, t, nk)
	msg := &Test{
		AString: name,
		AInt:    amount,
	}
	res := new(Test)
	req := Rpc("protoTest", msg, res)
	if err := req.WithHttpKey(nk.Name()).Do(ctx, cl); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("AString: %s", res.AString)
	t.Logf("AInt: %d", res.AInt)
	if res.AString != "hello "+name {
		t.Errorf("expected %q, got: %q", "hello "+name, res.AString)
	}
	if res.AInt != 2*amount {
		t.Errorf("expected %d, got: %d", 2*amount, res.AInt)
	}
}

func TestAuthenticateDevice(t *testing.T) {
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	defer cancel()
	cl := newClient(ctx, t, nk, WithServerKey(nk.ServerKey()))
	createAccount(ctx, t, cl)
}

func TestPing(t *testing.T) {
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	defer cancel()
	cl := newClient(ctx, t, nk, WithServerKey(nk.ServerKey()))
	conn := createAccountAndConn(ctx, t, cl)
	defer conn.Close()
	if err := conn.Ping(ctx); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if len(conn.l) != 0 {
		t.Errorf("expected len(conn.l) == 0, got: %d", len(conn.l))
	}
	errc := make(chan error, 1)
	conn.PingAsync(ctx, func(err error) {
		defer close(errc)
		errc <- err
	})
	select {
	case <-ctx.Done():
	case err := <-errc:
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	}
	if len(conn.l) != 0 {
		t.Errorf("expected len(conn.l) == 0, got: %d", len(conn.l))
	}
}

func TestRpcRealtime(t *testing.T) {
}

func TestChannels(t *testing.T) {
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	defer cancel()
	cl := newClient(ctx, t, nk, WithServerKey(nk.ServerKey()))
	conn := createAccountAndConn(ctx, t, cl)
	defer conn.Close()
}

func newClient(ctx context.Context, t *testing.T, nk *nktest.Runner, opts ...Option) *Client {
	urlstr, err := nktest.RunProxy(ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("url: %s", urlstr)
	opts = append([]Option{
		WithURL(urlstr),
		WithServerKey(nk.ServerKey()),
		WithTransport(&http.Transport{
			DisableCompression: true,
		}),
	}, opts...)
	return New(opts...)
}

func createAccount(ctx context.Context, t *testing.T, cl *Client) {
	deviceId := uuid.New().String()
	t.Logf("registering: %s", deviceId)
	if err := cl.AuthenticateDevice(ctx, deviceId, true, ""); err != nil {
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

func createAccountAndConn(ctx context.Context, t *testing.T, cl *Client, opts ...ConnOption) *Conn {
	createAccount(ctx, t, cl)
	conn, err := cl.NewConn(ctx, append([]ConnOption{WithConnFormat("json")}, opts...)...)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	return conn
}

type rewards struct {
	Rewards int64 `json:"rewards,omitempty"`
}
