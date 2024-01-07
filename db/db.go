package db

import (
	"context"
	"database/sql"
	"github.com/artela-network/galxe-integration/config"
	log "github.com/sirupsen/logrus"
	"strings"
)

func GetDB(ctx context.Context, config *config.DBConfig) (*sql.DB, string, error) {
	split := strings.Split(config.URL, "://")
	if len(split) != 2 {
		log.Fatalf("invalid db connection info: %s", config.URL)
	}
	driver := split[0]

	// Hardcode postgres here for now
	conn, err := newPostgres(ctx, config)
	return conn, driver, err
}
