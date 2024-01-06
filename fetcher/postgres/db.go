package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/artela-network/galxe-integration/fetcher"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
)

const driver = "postgres"

type postgresDAO struct {
	conn *sql.DB
}

func createIndex(conn *sql.DB, indexName, tableName, columnName string) {
	query := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s);", indexName, tableName, columnName)
	_, err := conn.Exec(query)
	if err != nil {
		log.Fatalf("Error creating index %s: %q", indexName, err)
	} else {
		log.Infof("Index %s created successfully.", indexName)
	}
}

func newPostgresDAO(ctx context.Context, dbConn string) fetcher.DAO {
	db, err := sql.Open(driver, dbConn)
	if err != nil {
		log.Fatal(err)
	}

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()

	if err := db.PingContext(timeoutCtx); err != nil {
		log.Fatal("db is not responding", err)
	}

	return &postgresDAO{
		conn: db,
	}
}

func (dao *postgresDAO) Init() fetcher.DAO {
	createTableSQL := `
        CREATE TABLE IF NOT EXISTS block_status (
            id SERIAL PRIMARY KEY,
            block_number INTEGER NOT NULL UNIQUE,
            status INTEGER NOT NULL,
            retry_count INTEGER DEFAULT 0,
            last_retry_at TIMESTAMP
        );`
	_, err := dao.conn.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	createIndex(dao.conn, "status_index", "block_status", "status")
	createIndex(dao.conn, "last_retry_at_index", "block_status", "last_retry_at")

	return dao
}

func (dao *postgresDAO) AddBlock(blockNumber uint64, status fetcher.BlockStatus) error {
	_, err := dao.conn.Exec("INSERT INTO block_status (block_number, status) VALUES ($1, $2) ON CONFLICT (block_number) DO NOTHING", blockNumber, status)
	return err
}

func (dao *postgresDAO) UpdateBlockStatus(blockNumber uint64, status fetcher.BlockStatus) error {
	_, err := dao.conn.Exec("UPDATE block_status SET status = $1 WHERE block_number = $2", status, blockNumber)
	return err
}

func (dao *postgresDAO) MigrateBlockStatus(blockNumber uint64, from fetcher.BlockStatus, to fetcher.BlockStatus) error {
	_, err := dao.conn.Exec("UPDATE block_status SET status = $1 WHERE block_number = $2 AND status = $3", to, blockNumber, from)
	return err
}

func (dao *postgresDAO) GetUnprocessedBlocks() ([]uint64, error) {
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

func (dao *postgresDAO) GetRetryBlocks(maxRetry uint64, retryThreshold time.Duration) ([]uint64, error) {
	rows, err := dao.conn.Query("SELECT block_number FROM block_status WHERE status = 3 AND retry_count < $1 AND (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) - EXTRACT(EPOCH FROM last_retry_at)) > $2", maxRetry, int64(retryThreshold.Seconds()))
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

func (dao *postgresDAO) MarkBlockForRetry(blockNumber uint64, maxRetry uint64) error {
	_, err := dao.conn.Exec("UPDATE block_status SET status = $1, retry_count = retry_count + 1, last_retry_at = CURRENT_TIMESTAMP WHERE block_number = $2 AND retry_count < $3",
		fetcher.StatusRetry, blockNumber, maxRetry)
	return err
}

func (dao *postgresDAO) GetLatestProcessedBlock() (uint64, error) {
	var latestBlock uint64
	row := dao.conn.QueryRow("SELECT block_number FROM block_status WHERE status = $1 ORDER BY block_number DESC LIMIT 1", fetcher.StatusProcessed)
	err := row.Scan(&latestBlock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return latestBlock, nil
}

func (dao *postgresDAO) GetBlockStatus(blockNumber uint64) (fetcher.BlockStatus, error) {
	var status fetcher.BlockStatus
	row := dao.conn.QueryRow("SELECT status FROM block_status WHERE block_number = $1", blockNumber)
	err := row.Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("block number %d not found", blockNumber)
		}
		return 0, err
	}
	return status, nil
}

func (dao *postgresDAO) ResetStaleProcessingBlocks(threshold time.Duration) error {
	_, err := dao.conn.Exec(
		"UPDATE block_status SET status = $1 WHERE status = $2 AND (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) - EXTRACT(EPOCH FROM last_retry_at)) > $3",
		fetcher.StatusUnprocessed, fetcher.StatusProcessing, int64(threshold.Seconds()))
	return err
}
