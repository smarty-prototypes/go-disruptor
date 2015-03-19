package disruptor

import (
	"sync/atomic"
	"time"
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

func (this *SharedWriter) Reserve(count int64) int64 {

	for {
		previous := this.written.Load() // we've written up to this point;
		upper := previous + count

		for upper-this.capacity > this.upstream.Read(0) {
			time.Sleep(time.Nanosecond)
		}

		if atomic.CompareAndSwapInt64(&this.written.sequence, previous, upper) {
			return upper
		}
	}
}

func (this *SharedWriter) Commit(lower, upper int64) {
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
