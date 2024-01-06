package noop

import (
	"github.com/artela-network/galxe-integration/indexer"
)

func init() {
	indexer.GetRegistry().Register(IndexerName, newNoopIndexer)
}
