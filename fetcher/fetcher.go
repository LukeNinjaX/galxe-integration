package fetcher

import (
	"context"
	"github.com/artela-network/galxe-integration/common"
	"github.com/artela-network/galxe-integration/config"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type fetcher struct {
	client              *ethclient.Client
	blockCache          chan *types.Block
	blockFetchTaskCache chan uint64
	dao                 *fetcherDAO
	ctx                 context.Context
	pullInterval        time.Duration
	retryInterval       time.Duration
	pollThread          uint64
	blockMaxRetry       uint64
	beginBlock          uint64
	maxProcessingTime   time.Duration

	indexers []common.Indexer
}

func NewFetcher(ctx context.Context, config *config.FetcherConfig) (common.Fetcher, error) {
	config.FillDefaults()

	rpcClient, err := rpc.DialContext(ctx, config.EthereumRPCUrl)
	if err != nil {
		log.Error("failed to dial ethereum rpc", err)
		return nil, err
	}

	client := ethclient.NewClient(rpcClient)
	maxProcessingTime, err := time.ParseDuration(config.MaxProcessingTime)
	if err != nil {
		log.Error("failed to parse max processing time", err)
		return nil, err
	}

	return &fetcher{
		ctx:                 ctx,
		client:              client,
		blockCache:          make(chan *types.Block, config.BlockCacheSize),
		blockFetchTaskCache: make(chan uint64, config.BlockCacheSize),
		dao:                 newFetcherDAO(ctx, config.DBConn).init(),
		pullInterval:        time.Duration(config.PullIntervalMs) * time.Millisecond,
		retryInterval:       time.Duration(config.RetryIntervalMs) * time.Millisecond,
		pollThread:          config.PollThread,
		blockMaxRetry:       config.BlockMaxRetry,
		beginBlock:          config.BeginBlock,
		maxProcessingTime:   maxProcessingTime,
	}, nil
}

func (f *fetcher) RegisterIndexer(indexer common.Indexer) {
	if f.indexers == nil {
		f.indexers = make([]common.Indexer, 0, 1)
	}
	f.indexers = append(f.indexers, indexer)
}

func (f *fetcher) Start() {
	for i := uint64(0); i < f.pollThread; i++ {
		go f.createWorker(i)
	}

	go f.createBlockListener()
	go f.createEventDispatcher()
	go f.monitorStaleProcessingTasks()
}

func (f *fetcher) createBlockListener() {
	lastPollTime := int64(0)
	for {
		select {
		case <-f.ctx.Done():
			log.Info("[block listener]: fetcher block listener stopped")
			return
		default:
			startTime := time.Now().UnixMilli()
			waitTime := startTime - lastPollTime
			log.Debugf("[block listener]: already waited %d ms", waitTime)
			if waitTime < f.pullInterval.Milliseconds() {
				remainingWaitTime := f.pullInterval.Milliseconds() - waitTime
				log.Debugf("[block listener]: still need to wait %d ms", remainingWaitTime)
				time.Sleep(time.Duration(remainingWaitTime) * time.Millisecond)
			}

			header, err := f.client.HeaderByNumber(f.ctx, nil)
			if err != nil {
				log.Error("[block listener]: error fetching latest block header:", err)
				time.Sleep(f.retryInterval)
				continue
			}

			lastPolledBlockHeight, err := f.dao.getLatestProcessedBlock()
			if err != nil {
				log.Error("[block listener]: failed to load latest processed block", err)
				return
			}
			lastPolledBlockHeight = max(lastPolledBlockHeight, f.beginBlock)

			currentBlock := header.Number.Uint64()
			for i := lastPolledBlockHeight + 1; i <= currentBlock; i++ {
				if err := f.dao.addBlock(i, StatusUnprocessed); err != nil {
					log.Error("[block listener]: failed to add block task", err)
					break
				}
			}

			unprocessedBlocks, err := f.dao.getUnprocessedBlocks(f.blockMaxRetry)
			if err != nil {
				log.Error("[block listener]: failed to load processed block", err)
			} else {
				for _, block := range unprocessedBlocks {
					log.Debugf("[block listener]: submitting block task %d", block)
					f.blockFetchTaskCache <- block
					log.Debugf("[block listener]: submitted block task %d", block)
				}
			}

			lastPollTime = time.Now().UnixMilli()
		}
	}
}

func (f *fetcher) createEventDispatcher() {
	for {
		select {
		case <-f.ctx.Done():
			log.Info("[event dispatcher]: stopped")
			return
		case block := <-f.blockCache:
			log.Debugf("[event dispatcher]: start dispatching block %d", block.NumberU64())
			if err := f.dao.migrateBlockStatus(block.NumberU64(), StatusUnprocessed, StatusProcessing); err != nil {
				log.Errorf("[event dispatcher]: failed to update block status to prcessing: %v", err)
				continue
			}
			receipts, err := f.client.BlockReceipts(f.ctx, rpc.BlockNumberOrHashWithHash(block.Hash(), true))
			if err != nil {
				log.Errorf("[event dispatcher]: error fetching block receipts: %v", err)
				continue
			}
			var processErr error
			for i, tx := range block.Transactions() {
				if tx.To() == nil {
					log.Debugf("[event dispatcher]: ignore contract creation tx %s", tx.Hash().Hex())
					continue
				}

				receipt := receipts[i]
				resChs := make([]chan error, 0, len(f.indexers))
				for _, indexer := range f.indexers {
					resCh := make(chan error, 1)
					eventCtx := &common.EventContext{
						BlockHeader: block.Header(),
						Transaction: tx,
						Receipt:     receipt,
						ResultChan:  resCh,
					}
					resChs = append(resChs, resCh)

					go func(indexer common.Indexer, eventCtx *common.EventContext) {
						log.Debugf("[event dispatcher]: submitting event task [block %d]->[tx %d]", block.NumberU64(), i)
						select {
						case <-f.ctx.Done():
							log.Info("[event dispatcher]: stopped")
						case indexer.Input() <- eventCtx:
							log.Debugf("[event dispatcher]: submitted event task [block %d]->[tx %d]", block.NumberU64(), i)
						}
					}(indexer, eventCtx)
				}

				for _, resCh := range resChs {
					select {
					case <-f.ctx.Done():
						log.Info("[event dispatcher]: stopped")
					case err := <-resCh:
						if err != nil {
							log.Errorf("[event dispatcher]: error dispatching event: %v", err)
							processErr = err
							break
						}
					}
				}
			}
			if processErr != nil {
				if err := f.dao.markBlockForRetry(block.NumberU64(), f.blockMaxRetry); err != nil {
					log.Errorf("[event dispatcher]: failed to mark block for retry: %v", err)
				}
			} else {
				if err := f.dao.updateBlockStatus(block.NumberU64(), StatusProcessed); err != nil {
					log.Errorf("[event dispatcher]: failed to update block status to processed: %v", err)
					continue
				}
				log.Infof("[event dispatcher]: processed block %d", block.NumberU64())
			}
		}
	}
}

func (f *fetcher) createWorker(index uint64) {
	for {
		select {
		case <-f.ctx.Done():
			log.Infof("[fetcher worker%d]: stopped", index)
			return
		case blockNum := <-f.blockFetchTaskCache:
			log.Debugf("[fetcher worker%d]: start fetching block %d", index, blockNum)

			status, err := f.dao.getBlockStatus(blockNum)
			if err != nil {
				log.Errorf("[fetcher worker%d]: failed to load block status", err)
				continue
			}

			if status != StatusUnprocessed {
				log.Debugf("[fetcher worker%d]: block %d is already processed or processing, skipping", index, blockNum)
				continue
			}

			block, err := f.client.BlockByNumber(f.ctx, new(big.Int).SetUint64(blockNum))
			if err != nil {
				log.Errorf("[fetcher worker%d]: error fetching block %d: %v", index, blockNum, err)
				continue
			}

			f.blockCache <- block
			log.Debugf("[fetcher worker%d]: fetched block %d", index, blockNum)
		}
	}
}

func (f *fetcher) monitorStaleProcessingTasks() {
	// checks every 1min for long processing tasks
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-f.ctx.Done():
			log.Info("[stale processing task monitor]: stopped")
			return
		case <-ticker.C:
			log.Info("[stale processing task monitor]: checking for stale processing tasks")
			if err := f.dao.resetStaleProcessingBlocks(f.maxProcessingTime); err != nil {
				log.Errorf("[stale processing task monitor]: failed to reset stale processing blocks: %v", err)
			}
		}
	}
}
