package router

import (
	"errors"
	"sync"

	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/messages"
	"github.com/mitghi/protox/protobase"
)

// Logger is the default, structured package level logging service.
var (
	logger protobase.LoggingInterface
)

func init() {
	logger = logging.NewLogger("Router")
}

// Error messages
var (
	// RRemoveErr indicates a problem associated wiht removing
	// a path from router.
	RRemoveErr error = errors.New("router: cannot remove")
	// SNoSubs indicates that topic does not exist.
	SNoSubs error = errors.New("subs: cannot locate subscriptions.")
)

type Router struct {
	sync.RWMutex

	subs   *Subs
	retain *messages.Retain
}
