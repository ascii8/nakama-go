package main

import (
	"bytes"
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
	logger.Info("INIT")
	if err := initializer.RegisterRpc("dailyRewards", dailyRewards); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("protoTest", protoTest); err != nil {
		return err
	}
	return nil
}

func dailyRewards(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Info("dailyRewards payload: %q", payload)
	// decode request
	dec := json.NewDecoder(strings.NewReader(payload))
	dec.DisallowUnknownFields()
	var req Rewards
	if err := dec.Decode(&req); err != nil {
		return "", err
	}
	// encode response
	res := new(bytes.Buffer)
	if err := json.NewEncoder(res).Encode(Rewards{
		Rewards: req.Rewards * 2,
	}); err != nil {
		return "", err
	}
	return res.String(), nil
}

type Rewards struct {
	Rewards int64 `json:"rewards,omitempty"`
}

func protoTest(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Info("protoTest payload: %q", payload)
	req := new(testpb.Test)
	if err := protojson.Unmarshal([]byte(payload), req); err != nil {
		return "", fmt.Errorf("unable to unmarshal protobuf message: %w", err)
	}
	res := &testpb.Test{
		AString: "hello " + req.AString,
		AInt:    2 * req.AInt,
	}
	buf, err := protojson.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("unable to marshal protobuf message: %w", err)
	}
	return string(buf), nil
}
