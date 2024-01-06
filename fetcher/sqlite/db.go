package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/artela-network/galxe-integration/fetcher"
	log "github.com/sirupsen/logrus"
	"time"
)

const driver = "sqlite3"

type sqliteDAO struct {
	conn *sql.DB
}

func newSqliteDAO(ctx context.Context, dbConn string) fetcher.DAO {
	db, err := sql.Open(driver, dbConn)
	if err != nil {
		log.Fatal(err)
	}

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()

	if err := db.PingContext(timeoutCtx); err != nil {
		log.Fatal("db is not responding", err)
	}

	return &sqliteDAO{
		conn: db,
	}
}

func (dao *sqliteDAO) Init() fetcher.DAO {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS block_status (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			block_number INTEGER NOT NULL UNIQUE,
			status INTEGER NOT NULL,
			retry_count INTEGER DEFAULT 0,
			last_retry_at DATETIME
		);`
	_, err := dao.conn.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	return dao
}

func (dao *sqliteDAO) AddBlock(blockNumber uint64, status fetcher.BlockStatus) error {
	_, err := dao.conn.Exec("INSERT OR IGNORE INTO block_status (block_number, status) VALUES (?, ?)", blockNumber, status)
	return err
}

func (dao *sqliteDAO) UpdateBlockStatus(blockNumber uint64, status fetcher.BlockStatus) error {
	_, err := dao.conn.Exec("UPDATE block_status SET status = ? WHERE block_number = ?", status, blockNumber)
	return err
}

func (dao *sqliteDAO) MigrateBlockStatus(blockNumber uint64, from fetcher.BlockStatus, to fetcher.BlockStatus) error {
	_, err := dao.conn.Exec("UPDATE block_status SET status = ? WHERE block_number = ? AND status = ?", to, blockNumber, from)
	return err
}

func (dao *sqliteDAO) GetUnprocessedBlocks() ([]uint64, error) {
	rows, err := dao.conn.Query("SELECT block_number FROM block_status WHERE status = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blockNumbers []uint64
	for rows.Next() {
		var blockNumber uint64
		if err := rows.Scan(&blockNumber); err != nil {
			return nil, err
		}
		blockNumbers = append(blockNumbers, blockNumber)
	}

	return blockNumbers, nil
}

func (dao *sqliteDAO) GetRetryBlocks(maxRetry uint64, retryThreshold time.Duration) ([]uint64, error) {
	rows, err := dao.conn.Query("SELECT block_number FROM block_status WHERE status = 3 AND retry_count < ? AND (strftime('%s', 'now') - strftime('%s', last_retry_at)) > ?", maxRetry, retryThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blockNumbers []uint64
	for rows.Next() {
		var blockNumber uint64
		if err := rows.Scan(&blockNumber); err != nil {
			return nil, err
		}
		blockNumbers = append(blockNumbers, blockNumber)
	}

	return blockNumbers, nil
}

func (dao *sqliteDAO) MarkBlockForRetry(blockNumber uint64, maxRetry uint64) error {
	_, err := dao.conn.Exec("UPDATE block_status SET status = ?, retry_count = retry_count + 1, last_retry_at = CURRENT_TIMESTAMP WHERE block_number = ? AND retry_count < ?",
		fetcher.StatusUnprocessed, blockNumber, maxRetry)
	return err
}

func (dao *sqliteDAO) GetLatestProcessedBlock() (uint64, error) {
	var latestBlock uint64
	row := dao.conn.QueryRow("SELECT block_number FROM block_status WHERE status = ? ORDER BY block_number DESC LIMIT 1", fetcher.StatusProcessed)
	err := row.Scan(&latestBlock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return latestBlock, nil
}

func (dao *sqliteDAO) GetBlockStatus(blockNumber uint64) (fetcher.BlockStatus, error) {
	var status fetcher.BlockStatus
	row := dao.conn.QueryRow("SELECT status FROM block_status WHERE block_number = ?", blockNumber)
	err := row.Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("block number %d not found", blockNumber)
		}
		return 0, err
	}
	return status, nil
}

func (dao *sqliteDAO) ResetStaleProcessingBlocks(threshold time.Duration) error {
	_, err := dao.conn.Exec(
		"UPDATE block_status SET status = ? WHERE status = ? AND (strftime('%s', 'now') - strftime('%s', last_retry_at)) > ?",
		fetcher.StatusUnprocessed, fetcher.StatusProcessing, int64(threshold.Seconds()))
	return err
}
