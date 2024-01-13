package api

import (
	"context"
	"database/sql"
	"errors"
	"github.com/artela-network/galxe-integration/common"
	"github.com/artela-network/galxe-integration/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	router *gin.Engine
	server *http.Server
	db     *sql.DB

	ctx  context.Context
	conf *config.APIConfig

	fetcher  common.Fetcher
	indexers []common.Indexer
}

func NewServer(ctx context.Context, config *config.APIConfig, _ string, db *sql.DB, fetcher common.Fetcher, indexers []common.Indexer) *Server {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(log.StandardLogger().Out)
	gin.DefaultErrorWriter = io.MultiWriter(log.StandardLogger().Out)

	r := gin.Default()

	if config.Host == "" {
		config.Host = "127.0.0.1"
	}
	if config.Port == 0 {
		config.Port = 9211
	}

	// CORS for https://galxe.com and Setup CORS to allow specific origins, methods, and headers
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://galxe.com"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.RecoveryWithWriter(log.StandardLogger().Out))

	s := &Server{
		router: r,
		conf:   config,
		ctx:    ctx,
		server: &http.Server{
			Addr:    config.Host + ":" + strconv.Itoa(int(config.Port)),
			Handler: r,
		},
		db:       db,
		fetcher:  fetcher,
		indexers: indexers,
	}

	apiGroup := r.Group("/api")
	apiGroup.GET("/ping", s.ping)
	apiGroup.GET("/jit-gaming/:address", s.completedJITGaming)
	apiGroup.GET("/metrics", s.metrics)

	return s
}

func (s *Server) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (s *Server) metrics(c *gin.Context) {
	indexerMetrics := make(map[string]interface{})
	for _, indexer := range s.indexers {
		indexerMetrics[indexer.Name()] = indexer.Metrics()
	}
	fetcherMetrics := s.fetcher.Metrics()

	c.JSON(200, gin.H{
		"fetcher": fetcherMetrics,
		"indexer": indexerMetrics,
	})
}

func (s *Server) completedJITGaming(c *gin.Context) {
	ethAddress := c.Param("address")
	if strings.HasPrefix(ethAddress, "/") {
		ethAddress = ethAddress[1:]
	}
	if strings.HasSuffix(ethAddress, "/") {
		ethAddress = ethAddress[:len(ethAddress)-1]
	}

	if ethAddress == "" || len(ethAddress) != 42 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing Ethereum address",
		})
		return
	}

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM scored_players WHERE LOWER(player) = LOWER($1))", ethAddress).Scan(&exists)
	if err != nil {
		log.Errorf("Failed to query database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check address",
		})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"completed": exists,
	})
}

func (s *Server) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("start server fail: %v", err)
		}
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("shutdown server fail: %v", err)
	} else {
		log.Info("api server has been shutdown")
	}
}
