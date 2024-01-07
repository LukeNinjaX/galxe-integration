package db

import (
	"context"
	"database/sql"
	"github.com/artela-network/galxe-integration/config"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func newPostgres(ctx context.Context, dbConfig *config.DBConfig) (*sql.DB, error) {
	split := strings.Split(dbConfig.URL, "://")
	if len(split) != 2 {
		log.Fatalf("invalid db connection info: %s", dbConfig.URL)
	}
	conn := split[1]

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if dbConfig.MaxConnection == 0 {
		dbConfig.MaxConnection = 50
	}
	db.SetMaxOpenConns(int(dbConfig.MaxConnection))

	timeoutCtx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()

	if err := db.PingContext(timeoutCtx); err != nil {
		return nil, err
	}

	return db, nil
}
