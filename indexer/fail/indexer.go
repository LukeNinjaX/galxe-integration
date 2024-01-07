package fail

import (
	"context"
	"database/sql"
	"errors"
	"github.com/artela-network/galxe-integration/common"
	"github.com/artela-network/galxe-integration/config"
	log "github.com/sirupsen/logrus"
)

const IndexerName = "Fail"

func newFailIndexer(ctx context.Context, _ *config.IndexerConfig, _ string, _ *sql.DB) (common.Indexer, error) {
	indexer := &failIndexer{
		inputCh: make(chan *common.EventContext, 100),
		ctx:     ctx,
	}
	indexer.Run()

	return indexer, nil
}

type failIndexer struct {
	inputCh chan *common.EventContext
	ctx     context.Context
}

func (n *failIndexer) Input() chan<- *common.EventContext {
	return n.inputCh
}

func (n *failIndexer) Run() {
	go func() {
		for {
			select {
			case eventCtx := <-n.inputCh:
				log.Infof("[fail indexer] received new tx[%s] @ block[%d] ",
					eventCtx.Transaction.Hash().Hex(), eventCtx.BlockHeader.Number.Uint64())
				select {
				case <-n.ctx.Done():
					log.Info("[fail indexer] stopped")
					return
				case eventCtx.ResultChan <- errors.New("error"):
					log.Infof("[fail indexer] processed tx[%s] @ block[%d] ",
						eventCtx.Transaction.Hash().Hex(), eventCtx.BlockHeader.Number.Uint64())
				default:
					close(eventCtx.ResultChan)
					log.Info("[fail indexer] result chan full")
				}
			case <-n.ctx.Done():
				log.Info("[fail indexer] stopped")
				return
			}
		}
	}()
}
