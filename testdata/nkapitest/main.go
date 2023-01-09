package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	testpb "github.com/ascii8/nakama-go/testdata/proto"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	if err := initializer.RegisterRpc("dailyRewards", dailyRewards); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("protoTest", protoTest); err != nil {
		return err
	}
	return nil
}

func dailyRewards(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payloadstr string) (string, error) {
	// decode request
	dec := json.NewDecoder(strings.NewReader(payloadstr))
	dec.DisallowUnknownFields()
	var req Rewards
	if err := dec.Decode(&req); err != nil {
		return "", err
	}
	logger.WithField("req", req).Debug("dailyRewards")
	res := Rewards{
		Rewards: req.Rewards * 2,
	}
	logger.WithField("res", res).Debug("dailyRewards")
	// encode response
	buf, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

type Rewards struct {
	Rewards int64 `json:"rewards,omitempty"`
}

func protoTest(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payloadstr string) (string, error) {
	req := new(testpb.Test)
	if err := protojson.Unmarshal([]byte(payloadstr), req); err != nil {
		return "", fmt.Errorf("unable to unmarshal protobuf message: %w", err)
	}
	logger.WithField("req", req).Debug("protoTest")
	res := &testpb.Test{
		AString: "hello " + req.AString,
		AInt:    2 * req.AInt,
	}
	logger.WithField("res", res).Debug("protoTest")
	buf, err := protojson.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
