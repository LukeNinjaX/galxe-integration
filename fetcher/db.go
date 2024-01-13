package fetcher

import (
	"context"
	"database/sql"
	"sync"
	"time"
)

type BlockStatus int

const (
	StatusUnprocessed BlockStatus = iota
	StatusProcessing
	StatusProcessed
	StatusRetry
)

type DAO interface {
	Init() DAO
	AddBlock(blockNumber uint64, status BlockStatus) error
	UpdateBlockStatus(blockNumber uint64, status BlockStatus) error
	MigrateBlockStatus(blockNumber uint64, from BlockStatus, to BlockStatus) error
	GetUnprocessedBlocks() ([]uint64, error)
	GetRetryBlocks(maxRetry uint64, retryThreshold time.Duration) ([]uint64, error)
	MarkBlockForRetry(blockNumber uint64, maxRetry uint64) error
	GetLatestProcessedBlock() (uint64, error)
	GetBlockStatus(blockNumber uint64) (BlockStatus, error)
	ResetStaleProcessingBlocks(threshold time.Duration) error
	GetCountByBlockStatus(status BlockStatus) (uint64, error)
	GetMaxProcessedBlockNumber() (uint64, error)
}

type Builder func(ctx context.Context, db *sql.DB) DAO

var registry Registry

type Registry struct {
	daos sync.Map
}

func (r *Registry) Register(tpy string, builder Builder) {
	r.daos.Store(tpy, builder)
}

func (r *Registry) GetDAO(ctx context.Context, driver string, db *sql.DB) DAO {
	// Parsing the connection string (assuming it's in PostgreSQL format)
	builder, exist := r.daos.Load(driver)
	if !exist {
		return nil
	}

	return builder.(Builder)(ctx, db)
}

func GetRegistry() *Registry {
	return &registry
}
