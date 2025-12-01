package disruptor

import "math"

type multiSequencer struct {
	written   atomicSequence
	gate      atomicSequence
	upstream  sequenceBarrier
	capacity  int64
	mask      int64
	shift     uint8
	committed []int32
}

func newMultiSequencer(write *multiSequencerBarrier, upstream sequenceBarrier) *multiSequencer {
	return &multiSequencer{
		written:   write.written,
		upstream:  upstream,
		capacity:  write.capacity,
		gate:      newSequence(),
		mask:      write.mask,
		shift:     write.shift,
		committed: write.committed,
	}
}

func (this *multiSequencer) Reserve(count int64) int64 {
	for {
		previous := this.written.Load()
		upper := previous + count

		for upper-this.capacity > this.gate.Load() {
			this.gate.Store(this.upstream.Load())
		}

		if this.written.CompareAndSwap(previous, upper) {
			return upper
		}
	}
}

func (this *multiSequencer) Commit(lower, upper int64) {
	if lower == upper {
		this.committed[upper&this.mask] = int32(upper >> this.shift)
	} else {
		// working down the array keeps all items in the commit together
		// otherwise the reader(s) could split up the group
		for upper >= lower {
			this.committed[upper&this.mask] = int32(upper >> this.shift)
			upper--
		}
	}
}

type multiSequencerBarrier struct {
	written   atomicSequence
	committed []int32
	capacity  int64
	mask      int64
	shift     uint8
}

func newSharedWriterBarrier(written atomicSequence, capacity int64) *multiSequencerBarrier {
	committed := make([]int32, capacity)
	for i := range committed {
		committed[i] = int32(defaultSequenceValue)
	}

	return &multiSequencerBarrier{
		written:   written,
		committed: committed,
		capacity:  capacity,
		mask:      capacity - 1,
		shift:     uint8(math.Log2(float64(capacity))),
	}
}

func (this *multiSequencerBarrier) Load(lower int64) int64 {
	shift, mask := this.shift, this.mask
	upper := this.written.Load()

	for ; lower <= upper; lower++ {
		if this.committed[lower&mask] != int32(lower>>shift) {
			return lower - 1
		}
	}

	return upper
}
