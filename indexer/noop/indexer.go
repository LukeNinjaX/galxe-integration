package noop

import (
	"context"
	"github.com/artela-network/galxe-integration/common"
	"github.com/artela-network/galxe-integration/config"
	log "github.com/sirupsen/logrus"
)

const IndexerName = "Noop"

func newNoopIndexer(ctx context.Context, _ *config.IndexerConfig) (common.Indexer, error) {
	indexer := &noopIndexer{
		inputCh: make(chan *common.EventContext, 100),
		ctx:     ctx,
	}
	indexer.Run()

	return indexer, nil
}

type noopIndexer struct {
	inputCh chan *common.EventContext
	ctx     context.Context
}

func (n *noopIndexer) Input() chan<- *common.EventContext {
	return n.inputCh
}

func (n *noopIndexer) Run() {
	go func() {
		for {
			select {
			case eventCtx := <-n.inputCh:
				log.Infof("[noop indexer] received new tx[%s] @ block[%d] ",
					eventCtx.Transaction.Hash().Hex(), eventCtx.BlockHeader.Number.Uint64())
				select {
				case <-n.ctx.Done():
					log.Info("[noop indexer] stopped")
					return
				case eventCtx.ResultChan <- nil:
					log.Infof("[noop indexer] processed tx[%s] @ block[%d] ",
						eventCtx.Transaction.Hash().Hex(), eventCtx.BlockHeader.Number.Uint64())
				default:
					close(eventCtx.ResultChan)
					log.Info("result chan full")
				}
			case <-n.ctx.Done():
				log.Info("[noop indexer] stopped")
				return
			}
		}
	}()
}
