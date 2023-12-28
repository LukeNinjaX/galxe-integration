package api

import (
	"encoding/json"
	"github.com/artela-network/galxe-integration/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"sync"
)

type RouteRegistration struct {
	Method  string
	Handler gin.HandlerFunc
}

type Builder func(rawConf json.RawMessage) []RouteRegistration

var registry Registry

type Registry struct {
	handlers sync.Map
}

func (r *Registry) RegisterHandlerBuilder(handlerType string, builder Builder) {
	r.handlers.Store(handlerType, builder)
}

func (r *Registry) RegisterRoutes(router *gin.RouterGroup, rawConf json.RawMessage) {
	typeConf := &config.TypeConf{}
	if err := json.Unmarshal(rawConf, typeConf); err != nil {
		log.Panic("load bot handler config fail", err)
	}

	// iterate over sync.Map
	builder, exist := r.handlers.Load(typeConf.Type)
	if !exist {
		log.Panicf("bot handler type %s not found", typeConf.Type)
	}

	routes := builder.(Builder)(rawConf)
	for _, route := range routes {
		router.Handle(route.Method, "/"+typeConf.Type, route.Handler)
	}
}

func GetRegistry() *Registry {
	return &registry
}
