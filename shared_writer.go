package disruptor

import (
	"runtime"
	"sync/atomic"
)

type SharedWriter struct {
	written   *Cursor
	upstream  Barrier
	capacity  int64
	mask      int64
	shift     uint8
	committed []int32
}

func NewSharedWriter(write *SharedWriterBarrier, upstream Barrier) *SharedWriter {
	return &SharedWriter{
		written:   write.written,
		upstream:  upstream,
		capacity:  write.capacity,
		mask:      write.mask,
		shift:     write.shift,
		committed: write.committed,
	}
}

func (this *SharedWriter) Reserve(id string, count int64) int64 {

	for {
		current := this.written.Load() // we've written up to this point;
		next := current + count

		for spin := int64(0); next-this.capacity > this.upstream.Read(0); spin++ {
			if spin&SpinMask == 0 {
				runtime.Gosched()
			}
		}

		if atomic.CompareAndSwapInt64(&this.written.sequence, current, next) {
			return next
		}
	}
}

func (this *SharedWriter) Commit(id string, lower, upper int64) {
	if lower == upper {
		this.committed[upper&this.mask] = int32(upper >> this.shift)
	} else {
		// working down the array rather than up keeps all items in the commit together
		// otherwise the reader(s) could split up the group
		for upper >= lower {
			this.committed[upper&this.mask] = int32(upper >> this.shift)
			upper--
		}
	}
}
