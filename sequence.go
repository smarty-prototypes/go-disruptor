package disruptor

import "sync/atomic"

// atomicSequence is a cache-line-padded atomic int64 used to track sequence positions without false sharing. The 56
// bytes of padding on each side ensure the embedded atomic.Int64 occupies its own cache line.
type atomicSequence struct {
	_ [7]int64 // 56B left padding
	atomic.Int64
	_ [7]int64 // 56B right padding
}

func newSequence() *atomicSequence {
	this := &atomicSequence{}
	this.Store(defaultSequenceValue)
	return this
}

// newSequences allocates a slice of *atomicSequence in a contiguous space in memory
func newSequences(count int) []*atomicSequence {
	backing := make([]atomicSequence, count)
	sequences := make([]*atomicSequence, count)
	for i := range sequences {
		backing[i].Store(defaultSequenceValue)
		sequences[i] = &backing[i]
	}
	return sequences
}

const defaultSequenceValue = -1
