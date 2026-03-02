package disruptor

import "sync/atomic"

// atomicSequence is a cache-line-padded atomic int64 used to track sequence positions without false sharing. The 56
// bytes of padding on each side ensure the embedded atomic.Int64 occupies its own cache line.
type atomicSequence struct {
	_ [CacheLineBytes - 8]byte // 64B - 8B left padding
	atomic.Int64
	_ [CacheLineBytes - 8]byte // 64B - 8B right padding
}

func newSequence() *atomicSequence {
	this := &atomicSequence{}
	this.Store(defaultSequenceValue)
	return this
}

// newSequences allocates a slice of *atomicSequence in a contiguous space in memory
func newSequences(count int) []*atomicSequence {
	actual := make([]atomicSequence, count)
	this := make([]*atomicSequence, count)
	for i := range this {
		actual[i].Store(defaultSequenceValue)
		this[i] = &actual[i]
	}
	return this
}

const defaultSequenceValue = -1
