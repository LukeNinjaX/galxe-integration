package faucet

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/artela-network/galxe-integration/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	s, err := NewFaucet(nil, &config.FaucetConfig{URL: "http://47.251.61.27:8545", KeyFile: "../../privateKey.txt"})
	require.Equal(t, nil, err)
	defer s.client.Close()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		fmt.Println("sending transfer", i)
		go func() {
			hash, err := s.client.Transfer(s.privateKey, common.HexToAddress("0x22b7926DA60F97c1aAD776084174218EbBF05E28"), TransferAmount, s.getNonce())
			require.Equal(t, nil, err)
			_ = hash
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(100 * time.Second)
}
