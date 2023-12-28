package slack

import (
	"github.com/artela-network/galxe-integration/notifier"
)

func init() {
	notifier.GetRegistry().Register(NotifierName, NewSlackNotifier)
}
