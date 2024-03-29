package nakama

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ascii8/nktest"
	"github.com/google/uuid"
)

// TestMain handles setting up and tearing down the postgres and nakama
// containers.
func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx = nktest.WithAlwaysPullFromEnv(ctx, "PULL")
	ctx = nktest.WithUnderCIFromEnv(ctx, "CI")
	ctx = nktest.WithHostPortMap(ctx)
	var opts []nktest.BuildConfigOption
	if os.Getenv("CI") == "" {
		opts = append(opts, nktest.WithDefaultGoEnv(), nktest.WithDefaultGoVolumes())
	}
	nktest.Main(ctx, m,
		nktest.WithDir("."),
		nktest.WithBuildConfig("./nkapitest", opts...),
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
	conn := createAccountAndConn(ctx, t, cl, true)
	defer conn.Close()
	if err := conn.Ping(ctx); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if len(conn.m) != 0 {
		t.Errorf("expected len(conn.l) == 0, got: %d", len(conn.m))
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
	if len(conn.m) != 0 {
		t.Errorf("expected len(conn.l) == 0, got: %d", len(conn.m))
	}
}

func TestMatch(t *testing.T) {
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	defer cancel()
	cl1 := newClient(ctx, t, nk)
	conn1 := createAccountAndConn(ctx, t, cl1, true)
	defer conn1.Close()
	a1, err := cl1.Account(ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("account1: %+v", a1)
	joinCh := make(chan *MatchPresenceEventMsg, 1)
	defer close(joinCh)
	conn1.MatchPresenceEventHandler = func(_ context.Context, msg *MatchPresenceEventMsg) {
		joinCh <- msg
	}
	m1, err := conn1.MatchCreate(ctx, "")
	switch {
	case err != nil:
		t.Fatalf("expected no error, got: %v", err)
	case m1.MatchId == "":
		t.Fatalf("expected non-empty m1.MatchId")
	case m1.Authoritative:
		t.Errorf("expected m1.Authoritative == false")
	case m1.Size == 0:
		t.Errorf("expected m1.Size != 0")
	case m1.Self.UserId != a1.User.Id:
		t.Errorf("expected m1.Self.UserId == a1.User.Id")
	}
	for _, p := range m1.Presences {
		t.Logf("p %s: %v", p.UserId, p.Status)
	}
	cl2 := newClient(ctx, t, nk)
	conn2 := createAccountAndConn(ctx, t, cl2, true)
	defer conn2.Close()
	dataCh := make(chan *MatchDataMsg, 1)
	defer close(dataCh)
	conn2.MatchDataHandler = func(_ context.Context, msg *MatchDataMsg) {
		dataCh <- msg
	}
	m2, err := conn2.MatchJoin(ctx, m1.MatchId, nil)
	switch {
	case err != nil:
		t.Fatalf("expected no error, got: %v", err)
	case m2.MatchId == "":
		t.Fatalf("expected non-empty m2.MatchId")
	case m1.MatchId != m2.MatchId:
		t.Errorf("expected m1.MatchId == m2.MatchId")
	case m2.Authoritative:
		t.Errorf("expected m2.Authoritative == false")
	}
	a2, err := cl2.Account(ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("account2: %+v", a2)
	select {
	case <-ctx.Done():
		t.Fatalf("context closed: %v", ctx.Err())
	case msg := <-joinCh:
		switch {
		case len(msg.Joins) != 1:
			t.Fatalf("expected 1 join, got: %d", len(msg.Joins))
		case msg.Joins[0].UserId != a2.User.Id:
			// t.Logf("msg: %+v", msg)
			// t.Fatalf("expected msg.Joins[0].UserId (%s) == a2.User.Id (%s)", msg.Joins[0].UserId, a2.User.Id)
		}
	}
	if err := conn1.MatchDataSend(ctx, m1.MatchId, 1, []byte(`hello world`), true); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	select {
	case <-ctx.Done():
		t.Fatalf("context closed: %v", ctx.Err())
	case msg := <-dataCh:
		if s, exp := string(msg.Data), "hello world"; s != exp {
			t.Errorf("expected %q, got: %q", exp, s)
		}
	}
	if err := conn1.MatchLeave(ctx, m1.MatchId); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if err := conn2.MatchLeave(ctx, m2.MatchId); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRpcRealtime(t *testing.T) {
	// TODO
}

func TestChannels(t *testing.T) {
	const target = "my_room"
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	defer cancel()
	cl1 := newClient(ctx, t, nk)
	conn1 := createAccountAndConn(ctx, t, cl1, true)
	defer conn1.Close()
	cl2 := newClient(ctx, t, nk)
	conn2 := createAccountAndConn(ctx, t, cl2, true)
	defer conn2.Close()
	recv := make(chan *ChannelMessage)
	conn2.ChannelMessageHandler = func(ctx context.Context, msg *ChannelMessage) {
		recv <- msg
	}
	ch1, err := conn1.ChannelJoin(ctx, target, ChannelType_ROOM, true, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("created channel: %s", ch1.Id)
	if _, err := conn2.ChannelJoin(ctx, target, ChannelType_ROOM, true, false); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	msg := map[string]interface{}{
		"msg":  "hello",
		"code": float64(15),
	}
	if _, err := conn1.ChannelMessageSend(ctx, ch1.Id, msg); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	var recvMsg *ChannelMessage
	select {
	case <-ctx.Done():
		t.Errorf("did not receive message: %v", ctx.Err())
	case <-time.After(1 * time.Minute):
		t.Error("did not receive message: timeout hit")
	case recvMsg = <-recv:
		t.Logf("received: %v", msg)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(recvMsg.Content), &m); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !reflect.DeepEqual(m, msg) {
		t.Errorf("expected m == msg:\n%#v\ngot:\n%#v", msg, m)
	}
}

func TestPersist(t *testing.T) {
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	cl := newClient(ctx, t, nk)
	conn := createAccountAndConn(ctx, t, cl, false, WithConnPersist(true))
	<-time.After(2 * conn.backoffMin)
	if conn.stop == true {
		t.Errorf("expected conn.stop == false")
	}
	if conn.Connected() == false {
		t.Fatalf("expected conn.Connected() == true")
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("expected on error, got: %v", err)
	}
	if conn.stop == false {
		t.Errorf("expected conn.stop == true")
	}
	if conn.Connected() == true {
		t.Errorf("expected conn.Connected() == false")
	}
	connectCh := make(chan bool, 1)
	conn.ConnectHandler = func(context.Context) {
		t.Logf("connected")
		connectCh <- true
	}
	disconnectCh := make(chan error, 1)
	conn.DisconnectHandler = func(_ context.Context, err error) {
		t.Logf("disconnected: %v", err)
		disconnectCh <- err
	}
	if err := conn.Open(ctx); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	<-time.After(2 * conn.backoffMin)
	if conn.stop == true {
		t.Errorf("expected conn.stop == false")
	}
	if conn.Connected() == false {
		t.Fatalf("expected conn.Connected() == true")
	}
	select {
	case <-ctx.Done():
		t.Fatalf("expected no error, got: %v", ctx.Err())
	case <-time.After(4 * conn.backoffMin):
		t.Fatalf("expected a connect event within %v", 4*conn.backoffMin)
	case b := <-connectCh:
		if b == false {
			t.Errorf("expected true")
		}
		t.Logf("connected: %t", b)
	}
	if err := conn.CloseWithStopErr(false, false, nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if conn.stop == true {
		t.Errorf("expected conn.stop == false")
	}
	<-time.After(4 * conn.backoffMin)
	if conn.Connected() == false {
		t.Errorf("expected conn.Connected() == true")
	}
	select {
	case <-ctx.Done():
		t.Fatalf("expected no error, got: %v", ctx.Err())
	case err := <-disconnectCh:
		if err != nil {
			t.Logf("disconnected!")
			t.Fatalf("expected no error, got: %v", err)
		}
	case <-time.After(conn.backoffMax):
		t.Errorf("expected a disconnect event within %v", conn.backoffMax)
	}
	// check no disconnect event received
	select {
	case <-ctx.Done():
		t.Errorf("expected no error, got: %v", ctx.Err())
		return
	case err := <-disconnectCh:
		t.Errorf("expected no disconnect event, got: %v", err)
	case <-time.After(conn.backoffMax):
		t.Logf("no disconnect")
	}
	if err := conn.CloseWithStopErr(true, true, errors.New("STOPPING")); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if conn.stop == false {
		t.Errorf("expected conn.stop == true")
	}
	<-time.After(4 * conn.backoffMin)
	if conn.Connected() == true {
		t.Errorf("expected conn.Connected() == false")
	}
	select {
	case <-ctx.Done():
		t.Fatalf("expected no error, got: %v", ctx.Err())
	case err := <-disconnectCh:
		switch {
		case err == nil:
			t.Error("expected disconnect event!")
		case err.Error() != "STOPPING":
			t.Error("expected STOPPING error")
		}
	}
	switch {
	case conn.Connected() == true:
		t.Errorf("expceted conn.Connected() == false")
	case conn.stop == false:
		t.Errorf("expected conn.stop == true")
	}
	cancel()
	<-time.After(conn.backoffMax)
	/*
		close(connectCh)
		close(disconnectCh)
	*/
}

func TestKeep(t *testing.T) {
	keep := os.Getenv("KEEP")
	if keep == "" {
		return
	}
	d, err := time.ParseDuration(keep)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	ctx, cancel, nk := nktest.WithCancel(context.Background(), t)
	defer cancel()
	urlstr, err := nk.RunProxy(ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("local: %s", nk.HttpLocal())
	t.Logf("grpc: %s", nk.GrpcLocal())
	t.Logf("http: %s", nk.HttpLocal())
	t.Logf("console: %s", nk.ConsoleLocal())
	t.Logf("http_key: %s", nk.HttpKey())
	t.Logf("server_key: %s", nk.ServerKey())
	t.Logf("proxy: %s", urlstr)
	select {
	case <-time.After(d):
	case <-ctx.Done():
	}
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

func createAccountAndConn(ctx context.Context, t *testing.T, cl *Client, check bool, opts ...ConnOption) *Conn {
	createAccount(ctx, t, cl)
	conn, err := cl.NewConn(ctx, append([]ConnOption{WithConnFormat("json")}, opts...)...)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if check && conn.Connected() != true {
		t.Fatalf("expected conn.Connected() == true")
	}
	return conn
}

type rewards struct {
	Rewards int64 `json:"rewards,omitempty"`
}
