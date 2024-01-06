package fetcher

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

type BlockStatus int

const (
	StatusUnprocessed BlockStatus = iota
	StatusProcessing
	StatusProcessed
)

type DAO interface {
	Init() DAO
	AddBlock(blockNumber uint64, status BlockStatus) error
	UpdateBlockStatus(blockNumber uint64, status BlockStatus) error
	MigrateBlockStatus(blockNumber uint64, from BlockStatus, to BlockStatus) error
	GetUnprocessedBlocks(retryCount uint64) ([]uint64, error)
	MarkBlockForRetry(blockNumber uint64, maxRetry uint64) error
	GetLatestProcessedBlock() (uint64, error)
	GetBlockStatus(blockNumber uint64) (BlockStatus, error)
	ResetStaleProcessingBlocks(threshold time.Duration) error
}

type Builder func(ctx context.Context, dbConn string) DAO

var registry Registry

type Registry struct {
	daos sync.Map
}

func (r *Registry) Register(tpy string, builder Builder) {
	r.daos.Store(tpy, builder)
}

func (r *Registry) GetDAO(ctx context.Context, dbConn string) DAO {
	// Parsing the connection string (assuming it's in PostgreSQL format)
	split := strings.Split(dbConn, "://")
	if len(split) != 2 {
		log.Fatalf("invalid db connection info: %s", dbConn)
	}
	driver, conn := split[0], split[1]

	builder, exist := r.daos.Load(driver)
	if !exist {
		return nil
	}

	return builder.(Builder)(ctx, conn)
}

func GetRegistry() *Registry {
	return &registry
}
