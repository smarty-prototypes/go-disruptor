package disruptor

import "math"

type multiSequencer struct {
	written   atomicSequence
	gate      atomicSequence
	upstream  sequenceBarrier
	committed []int32
	capacity  int64
	shift     uint8
}
type multiSequencerBarrier struct {
	written   atomicSequence
	committed []int32
	capacity  int64
	shift     uint8
}

func (this *multiSequencer) Reserve(count int64) int64 {
	// blocks until desired number of slots becomes available
	for {
		previous := this.written.Load()
		upper := previous + count

		for upper-this.capacity > this.gate.Load() {
			this.gate.Store(this.upstream.Load())
		}

		if this.written.CompareAndSwap(previous, upper) {
			return upper
		}

		// TODO: should we pass context.Context into this? if the caller aborts, we can skip the reservation request
	}
}
func (this *multiSequencer) Commit(lower, upper int64) {
	for mask := this.capacity - 1; upper >= lower; {
		this.committed[upper&mask] = int32(upper >> this.shift)
		upper--
	}
}

func (this *multiSequencerBarrier) Load(lower int64) int64 {
	upper := this.written.Load()

	for mask := this.capacity - 1; lower <= upper; lower++ {
		if this.committed[lower&mask] != int32(lower>>this.shift) {
			return lower - 1
		}
	}

	return upper
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type multiSequencerConfiguration struct {
	written   atomicSequence
	committed []int32
	capacity  int64
	shift     uint8
}

func newMultiSequencerConfiguration(written atomicSequence, capacity uint32) *multiSequencerConfiguration {
	committed := make([]int32, capacity)
	for i := range committed {
		committed[i] = int32(defaultSequenceValue)
	}

	return &multiSequencerConfiguration{
		written:   written,
		committed: committed,
		capacity:  int64(capacity),
		shift:     uint8(math.Log2(float64(capacity))),
	}
}

func (this *multiSequencerConfiguration) NewBarrier() *multiSequencerBarrier {
	return &multiSequencerBarrier{
		written:   this.written,
		committed: this.committed,
		capacity:  this.capacity,
		shift:     this.shift,
	}
}
func (this *multiSequencerConfiguration) NewSequencer(upstream sequenceBarrier) Sequencer {
	return &multiSequencer{
		written:   this.written,
		gate:      newSequence(),
		upstream:  upstream,
		committed: this.committed,
		capacity:  this.capacity,
		shift:     this.shift,
	}
}
