package common

type Indexer interface {
	Input() chan<- *EventContext
}

type Fetcher interface {
	RegisterIndexer(indexer Indexer)
	Start()
}

type Notifier interface {
	Notify(msg, from string, throttle bool)
}
