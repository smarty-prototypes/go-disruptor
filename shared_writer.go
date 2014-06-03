package disruptor

import "sync/atomic"

type SharedWriter struct {
	written   *Cursor
	upstream  Barrier
	capacity  int64
	gate      *Cursor
	mask      int64
	shift     uint8
	committed []int32
}

func NewSharedWriter(write *SharedWriterBarrier, upstream Barrier) *SharedWriter {
	return &SharedWriter{
		written:   write.written,
		upstream:  upstream,
		capacity:  write.capacity,
		gate:      NewCursor(),
		mask:      write.mask,
		shift:     write.shift,
		committed: write.committed,
	}
}

func (this *SharedWriter) Reserve(count int64) int64 {
	for {
		previous := this.written.Load()
		upper := previous + count

		if previous-this.capacity > this.gate.Load() {
			this.gate.Store(this.upstream.Read(0))
		}

		if atomic.CompareAndSwapInt64(&this.written.sequence, previous, upper) {
			return upper
		}
	}
}

func (this *SharedWriter) Commit(lower, upper int64) {
	// POTENTIAL TODO: start from upper and work toward lower
	// this may have the effect of keeping a batch together which
	// might otherwise be split up...
	for shift, mask := this.shift, this.mask; lower <= upper; lower++ {
		this.committed[lower&mask] = int32(lower >> shift)
	}
}
