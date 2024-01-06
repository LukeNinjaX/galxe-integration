package fail

import (
	"github.com/artela-network/galxe-integration/indexer"
)

func init() {
	indexer.GetRegistry().Register(IndexerName, newFailIndexer)
}
