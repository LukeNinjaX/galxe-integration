package common

type Measurable interface {
	Metrics() interface{}
}

type Indexer interface {
	Measurable
	Input() chan<- *EventContext
	Name() string
}

type Fetcher interface {
	Measurable
	RegisterIndexer(indexer Indexer)
	Start()
}

type Notifier interface {
	Notify(msg, from string, throttle bool)
}
