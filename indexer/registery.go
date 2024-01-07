package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artela-network/galxe-integration/common"
	"github.com/artela-network/galxe-integration/config"
	"sync"
)

type Builder func(ctx context.Context, config *config.IndexerConfig, driver string, db *sql.DB) (common.Indexer, error)

var registry Registry

type Registry struct {
	indexers sync.Map
}

func (r *Registry) Register(tpy string, builder Builder) {
	r.indexers.Store(tpy, builder)
}

func (r *Registry) GetIndexer(ctx context.Context, indexerConf *config.IndexerConfig, driver string, db *sql.DB) (common.Indexer, error) {
	builder, exist := r.indexers.Load(indexerConf.Type)
	if !exist {
		return nil, fmt.Errorf("indexer type %s not found", indexerConf.Type)
	}

	return builder.(Builder)(ctx, indexerConf, driver, db)
}

func GetRegistry() *Registry {
	return &registry
}
