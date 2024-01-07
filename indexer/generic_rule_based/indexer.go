package generic_rule_based

import (
	"context"
	"database/sql"
	"github.com/artela-network/galxe-integration/common"
	"github.com/artela-network/galxe-integration/config"
)

const IndexerName = "GenericRuleBased"

func newRuleBasedIndexer(_ context.Context, _ *config.IndexerConfig, _ string, _ *sql.DB) (common.Indexer, error) {
	panic("not implemented")
}
