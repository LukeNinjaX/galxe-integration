package cleaner

import (
	"database/sql"
	"time"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/onchain"
	log "github.com/sirupsen/logrus"
)

type cleaner struct {
	db *sql.DB
}

func NewCleaner(db *sql.DB) *cleaner {
	return &cleaner{db}
}

func (c *cleaner) Start() {
	go func() {
		for {
			rowsAffected, err := biz.LetTimeoutRecordRetry(c.db)
			if err != nil {
				log.Error("cleaner: reset status from 2 to 1 failed")
			}
			log.Debugf("cleaner: reset status from 2 to 1, %d rows affected", rowsAffected)
			time.Sleep(onchain.CleanDBInterval)
		}
	}()
}
