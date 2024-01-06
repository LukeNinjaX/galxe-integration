package sqlite

import "github.com/artela-network/galxe-integration/fetcher"

func init() {
	fetcher.GetRegistry().Register(driver, newSqliteDAO)
}
