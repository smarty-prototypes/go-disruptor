package disruptor

import "sync/atomic"

func NewCursor() *atomic.Int64 {
	this := &atomic.Int64{}
	this.Store(defaultCursorValue)
	return this
}

const defaultCursorValue = -1
