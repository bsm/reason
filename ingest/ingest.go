package ingest

import (
	"github.com/bsm/reason"
)

// DataStream defines a common interface for data streams.
type DataStream interface {
	// Next returns the next example. It will return io.EOF when there are no more examples.
	Next() (reason.Example, error)
}
