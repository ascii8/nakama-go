package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("INIT")
	if err := initializer.RegisterRpc("dailyRewards", dailyRewards); err != nil {
		return err
	}
	return nil
}

func dailyRewards(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Info("daily rewards payload: %q", payload)
	// decode request
	dec := json.NewDecoder(strings.NewReader(payload))
	dec.DisallowUnknownFields()
	var req Rewards
	if err := dec.Decode(&req); err != nil {
		return "", err
	}
	// encode response
	res := new(bytes.Buffer)
	enc := json.NewEncoder(res)
	if err := enc.Encode(Rewards{
		Rewards: req.Rewards * 2,
	}); err != nil {
		return "", err
	}
	return res.String(), nil
}

type Rewards struct {
	Rewards int64 `json:"rewards,omitempty"`
}
