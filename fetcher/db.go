package fetcher

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type BlockStatus int

const (
	StatusUnprocessed BlockStatus = iota
	StatusProcessing
	StatusProcessed
)

type fetcherDAO struct {
	conn *sql.DB
}

func newFetcherDAO(ctx context.Context, dbConn string) *fetcherDAO {
	split := strings.Split(dbConn, "://")
	if len(split) != 2 {
		log.Fatalf("invalid db connection info: %s", dbConn)
	}
	driver, conn := split[0], split[1]

	db, err := sql.Open(driver, conn)
	if err != nil {
		log.Fatal(err)
	}

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()

	if err := db.PingContext(timeoutCtx); err != nil {
		log.Fatal("db is not responding", err)
	}

	return &fetcherDAO{
		conn: db,
	}
}

func (dao *fetcherDAO) init() *fetcherDAO {
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

func (dao *fetcherDAO) addBlock(blockNumber uint64, status BlockStatus) error {
	_, err := dao.conn.Exec("INSERT OR IGNORE INTO block_status (block_number, status) VALUES (?, ?)", blockNumber, status)
	return err
}

func (dao *fetcherDAO) updateBlockStatus(blockNumber uint64, status BlockStatus) error {
	_, err := dao.conn.Exec("UPDATE block_status SET status = ? WHERE block_number = ?", status, blockNumber)
	return err
}

func (dao *fetcherDAO) migrateBlockStatus(blockNumber uint64, from BlockStatus, to BlockStatus) error {
	_, err := dao.conn.Exec("UPDATE block_status SET status = ? WHERE block_number = ? AND status = ?", to, blockNumber, from)
	return err
}

func (dao *fetcherDAO) getUnprocessedBlocks(retryCount uint64) ([]uint64, error) {
	rows, err := dao.conn.Query("SELECT block_number FROM block_status WHERE status = 0 and retry_count < ?", retryCount)
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

func (dao *fetcherDAO) markBlockForRetry(blockNumber uint64, maxRetry uint64) error {
	_, err := dao.conn.Exec("UPDATE block_status SET status = ?, retry_count = retry_count + 1, last_retry_at = CURRENT_TIMESTAMP WHERE block_number = ? AND retry_count < ?",
		StatusUnprocessed, blockNumber, maxRetry)
	return err
}

func (dao *fetcherDAO) getLatestProcessedBlock() (uint64, error) {
	var latestBlock uint64
	row := dao.conn.QueryRow("SELECT block_number FROM block_status WHERE status = ? ORDER BY block_number DESC LIMIT 1", StatusProcessed)
	err := row.Scan(&latestBlock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return latestBlock, nil
}

func (dao *fetcherDAO) getBlockStatus(blockNumber uint64) (BlockStatus, error) {
	var status BlockStatus
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

func (dao *fetcherDAO) resetStaleProcessingBlocks(threshold time.Duration) error {
	_, err := dao.conn.Exec(
		"UPDATE block_status SET status = ? WHERE status = ? AND (strftime('%s', 'now') - strftime('%s', last_retry_at)) > ?",
		StatusUnprocessed, StatusProcessing, int64(threshold.Seconds()))
	return err
}
