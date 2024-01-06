package main

import (
	// notifiers
	_ "github.com/artela-network/galxe-integration/notifier/slack"

	// indexers
	_ "github.com/artela-network/galxe-integration/indexer/generic_rule_based"

	// db
	_ "github.com/artela-network/galxe-integration/fetcher/postgres"
	_ "github.com/artela-network/galxe-integration/fetcher/sqlite"
)
