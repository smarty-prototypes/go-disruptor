package disruptor

// compositeBarrier is a barrier that returns the minimum sequence across multiple atomicSequence values. It
// represents a group of downstream consumers whose collective progress gates an upstream producer. Load iterates all
// sequences and returns the lowest value, which is the slowest consumer's position. The constructor optimizes for
// common cases: zero sequences returns an empty barrier, one sequence returns a plain atomicBarrier (no iteration).
type compositeBarrier []*atomicSequence

func newCompositeBarrier(sequences ...*atomicSequence) sequenceBarrier {
	if len(sequences) == 0 {
		return compositeBarrier{} // TODO: panic?
	} else if len(sequences) == 1 {
		return newAtomicBarrier(sequences[0])
	} else {
		return compositeBarrier(sequences)
	}
}

func (this compositeBarrier) Load(_ int64) int64 {
	var minimum int64 = 1<<63 - 1

	for _, item := range this {
		if sequence := item.Load(); sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}
