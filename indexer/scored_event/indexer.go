package scored_event

import (
	"context"
	"database/sql"
	"errors"
	"github.com/artela-network/galxe-integration/common"
	"github.com/artela-network/galxe-integration/config"
	"github.com/ethereum/go-ethereum/accounts/abi"
	eth "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
	"runtime"
	"strings"
	"time"
)

const IndexerName = "ScoredEvent"

var contractAddr = eth.HexToAddress("0xA8DF3c46212BDfA525a91F3c2FDb3C199281a60e")

type ScoredEvent struct {
	Player eth.Address
	Score  *big.Int
}

var scoredEventABI, _ = abi.JSON(strings.NewReader(`[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"player","type":"address"},{"indexed":false,"internalType":"uint256","name":"score","type":"uint256"}],"name":"Scored","type":"event"}]`))

func newScoredEventIndexer(ctx context.Context, conf *config.IndexerConfig, _ string, db *sql.DB) (common.Indexer, error) {
	// Create the scores table if it doesn't exist
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS scored_players (
        id SERIAL PRIMARY KEY,
        player VARCHAR(42) NOT NULL UNIQUE
    )`)
	if err != nil {
		log.Fatal("Failed to create scores table", err)
		return nil, err
	}

	if conf.Thread == 0 {
		conf.Thread = uint64(runtime.NumCPU())*2 + 1
	}

	indexer := &scoredEventIndexer{
		inputCh:     make(chan *common.EventContext, 100),
		ctx:         ctx,
		db:          db,
		concurrency: conf.Thread,
	}
	indexer.Run()

	return indexer, nil
}

type scoredEventIndexer struct {
	inputCh     chan *common.EventContext
	ctx         context.Context
	db          *sql.DB
	concurrency uint64
}

func (s *scoredEventIndexer) Input() chan<- *common.EventContext {
	return s.inputCh
}

func (s *scoredEventIndexer) Run() {
	go func() {
		for {
			log.Infof("[scored event indexer] currently there are %d tx waiting", len(s.inputCh))
			time.Sleep(5 * time.Second)
		}
	}()

	for i := uint64(0); i < s.concurrency; i++ {
		go func() {
			for {
				select {
				case eventCtx := <-s.inputCh:
					var err error
					for _, ethLog := range eventCtx.Receipt.Logs {
						// Check if the log's address matches the contract address
						if ethLog.Address != contractAddr {
							log.Debug("[scored event indexer] not target contract address")
							continue
						}

						scoredEventSig := scoredEventABI.Events["Scored"].ID
						if ethLog.Topics[0] != scoredEventSig {
							log.Debug("[scored event indexer] not scored event")
							continue
						}

						err = func() error {
							defer func() {
								if r := recover(); r != nil {
									log.Error("[scored event indexer] panic", r)
									err = errors.New("indexer panic")
								}
							}()
							event := new(ScoredEvent)
							if err := scoredEventABI.UnpackIntoInterface(event, "Scored", ethLog.Data); err != nil {
								log.Error("[scored event indexer] failed to unpack scored event", err)
								return err
							}

							log.Debugf("[scored event indexer] player %s scored %d", event.Player.Hex(), event.Score.Uint64())

							if (event.Player == eth.Address{}) {
								log.Debugf("[scored event indexer] npc player scored %d, ignore", event.Score.Uint64())
								return nil
							}

							if event.Score.Uint64() >= 5 {
								// we may receive duplicate logs here, need to ignore the conflicts
								_, err := s.db.Exec("INSERT INTO scored_players(player) VALUES($1) ON CONFLICT (player) DO NOTHING", event.Player.Hex())
								if err != nil {
									log.Error("[scored event indexer] failed to insert score", err)
									return err
								}
							} else {
								log.Debugf("[scored event indexer] player score %d is not 5", event.Score.Uint64())
							}
							return nil
						}()

						if err != nil {
							break
						}
					}

					select {
					case <-s.ctx.Done():
						log.Info("[scored event indexer] stopped")
						return
					case eventCtx.ResultChan <- err:
						log.Infof("[scored event indexer] processed tx[%s] @ block[%d] ",
							eventCtx.Transaction.Hash().Hex(), eventCtx.BlockHeader.Number.Uint64())
					default:
						close(eventCtx.ResultChan)
						log.Info("[scored event indexer] result chan full")
					}
				case <-s.ctx.Done():
					log.Info("[scored event indexer] stopped")
					return
				}
			}
		}()
	}
}
