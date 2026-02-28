package disruptor

// sequenceBarrier abstracts reading the committed or handled position from one or more sequences. Load returns the
// highest sequence that is safe to read up to from the given lower bound.
type sequenceBarrier interface {
	Load(int64) int64
}

// atomicBarrier is a sequenceBarrier backed by a single atomicSequence. Used when there is exactly one upstream
// sequence to track, avoiding the iteration overhead of compositeBarrier.
type atomicBarrier struct{ sequence *atomicSequence }

func newAtomicBarrier(sequence *atomicSequence) atomicBarrier {
	return atomicBarrier{sequence: sequence}
}
func (this atomicBarrier) Load(_ int64) int64 { return this.sequence.Load() }
