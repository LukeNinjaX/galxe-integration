package notifier

import (
	"context"
	"encoding/json"
	"github.com/artela-network/galxe-integration/common"
	"github.com/artela-network/galxe-integration/config"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Builder func(ctx context.Context, rawConf json.RawMessage) common.Notifier

var registry Registry

type Registry struct {
	notifiers sync.Map
}

func (r *Registry) Register(tpy string, builder Builder) {
	r.notifiers.Store(tpy, builder)
}

func (r *Registry) GetNotifier(ctx context.Context, rawConf json.RawMessage) common.Notifier {
	typeConf := &config.TypeConf{}
	if err := json.Unmarshal(rawConf, typeConf); err != nil {
		log.Panic("load config fail", err)
	}

	builder, exist := r.notifiers.Load(typeConf.Type)
	if !exist {
		return nil
	}

	return builder.(Builder)(ctx, rawConf)
}

func GetRegistry() *Registry {
	return &registry
}
