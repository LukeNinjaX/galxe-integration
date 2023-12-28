package common

import "github.com/ethereum/go-ethereum/core/types"

type EventContext struct {
	BlockHeader *types.Header
	Transaction *types.Transaction
	Receipt     *types.Receipt
	ResultChan  chan<- error
}
