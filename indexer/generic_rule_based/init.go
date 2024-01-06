package generic_rule_based

import (
	"github.com/artela-network/galxe-integration/indexer"
)

func init() {
	indexer.GetRegistry().Register(IndexerName, newRuleBasedIndexer)
}
