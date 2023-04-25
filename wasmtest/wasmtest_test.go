package wasmtest

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ascii8/nakama-go"
	"github.com/google/uuid"
)

func TestPersist(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl := newClient(t)
	conn := createAccountAndConn(ctx, t, cl, false, nakama.WithConnPersist(true))
	ch := make(chan error)
	conn.DisconnectHandler = func(_ context.Context, err error) {
		ch <- err
	}
	if err := conn.Open(ctx); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	<-time.After(1500 * time.Millisecond)
	if !conn.Connected() {
		t.Errorf("expected conn to be connected")
	}
	if err := conn.CloseWithStopErr(true, errors.New("STOPPING")); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	select {
	case <-ctx.Done():
		t.Fatalf("expected no error, got: %v", ctx.Err())
	case <-time.After(1 * time.Minute):
		t.Fatalf("expected disconnected error")
	case err := <-ch:
		switch {
		case err == nil:
			t.Fatalf("expected disconnected error")
		case err.Error() != "STOPPING":
			t.Errorf("expected STOPPING error")
		}
	}
	switch {
	case conn.Connected():
		t.Errorf("expected conn.Connected() == false")
	}
}

func newClient(t *testing.T) *nakama.Client {
	const urlstr = "http://127.0.0.1:7350"
	const serverKey = "nakama-go_server"
	t.Logf("url: %s", urlstr)
	opts := append([]nakama.Option{
		nakama.WithURL(urlstr),
		nakama.WithServerKey(serverKey),
		nakama.WithTransport(&http.Transport{
			DisableCompression: true,
		}),
	})
	return nakama.New(opts...)
}

func createAccount(ctx context.Context, t *testing.T, cl *nakama.Client) {
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
	found := false
	for _, d := range res.Devices {
		if d.Id == deviceId {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected accountRes.Devices to contain %s", deviceId)
	}
}

func createAccountAndConn(ctx context.Context, t *testing.T, cl *nakama.Client, check bool, opts ...nakama.ConnOption) *nakama.Conn {
	createAccount(ctx, t, cl)
	conn, err := cl.NewConn(ctx, append([]nakama.ConnOption{nakama.WithConnFormat("json")}, opts...)...)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if check && conn.Connected() != true {
		t.Fatalf("expected conn.Connected() == true")
	}
	return conn
}
