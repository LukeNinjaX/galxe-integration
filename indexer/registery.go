package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artela-network/galxe-integration/common"
	"github.com/artela-network/galxe-integration/config"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Builder func(ctx context.Context, config json.RawMessage) (common.Indexer, error)

var registry Registry

type Registry struct {
	indexers sync.Map
}

func (r *Registry) Register(tpy string, builder Builder) {
	r.indexers.Store(tpy, builder)
}

func (r *Registry) GetIndexer(ctx context.Context, indexerConf json.RawMessage) (common.Indexer, error) {
	var indexType string

	typeConf := &config.TypeConf{}
	if err := json.Unmarshal(indexerConf, typeConf); err != nil {
		log.Error("load config fail", err)
		return nil, err
	}

	builder, exist := r.indexers.Load(indexType)
	if !exist {
		return nil, fmt.Errorf("indexer type %s not found", indexType)
	}

	return builder.(Builder)(ctx, indexerConf)
}

func GetRegistry() *Registry {
	return &registry
}
