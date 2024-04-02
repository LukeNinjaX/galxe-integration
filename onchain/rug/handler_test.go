package rug

import (
	"fmt"
	"testing"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/config"
	"github.com/stretchr/testify/require"
)

func TestNewRug(t *testing.T) {
	cfg := &config.RugConfig{
		OnChain: config.OnChain{
			URL:     "https://betanet-inner3.artela.network",
			KeyFile: "../../rug.txt",
		},
	}
	s, err := NewRug(nil, cfg)
	require.Equal(t, nil, err)
	defer s.Client().Close()

	task := biz.AddressTask{
		ID:             123,
		GMTCreate:      time.Time{},
		GMTModify:      time.Time{},
		AccountAddress: new(string),
		TaskName:       new(string),
		TaskStatus:     new(string),
		Memo:           new(string),
		Txs:            new(string),
		TaskId:         new(string),
		TaskTopic:      new(string),
		JobBatchId:     new(string),
	}

	hash, err := s.send(task)
	require.Equal(t, nil, err)
	require.Equal(t, 1, len(hash))
	fmt.Println(hash[0].Hex())
}
