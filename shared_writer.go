package disruptor

import (
	"fmt"
	"runtime"
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

func (this *SharedWriter) Reserve(id string, count int64) int64 {

	for {
		// current := atomic.LoadInt64(&this.written.sequence)
		current := this.written.Load() // we've written up to this point;
		next := current + count

		for spin := int64(0); next-this.capacity > this.upstream.Read(0); spin++ {
			if spin&SpinMask == 0 {
				fmt.Printf("Writer '%s' gating on reader: Written up to Sequence: %d, Desired Sequence: %d, Waiting for: %d, Reader at: %d\n", id, current, next, (next - this.capacity), this.upstream.Read(0))
				time.Sleep(time.Second)
				runtime.Gosched()
			}
		}

		if atomic.CompareAndSwapInt64(&this.written.sequence, current, next) {
			fmt.Printf("Writer '%s' reserved up to sequence %d\n", id, next)
			return next
		} else {
			fmt.Printf("Writer '%s' collision on %d, trying again.\n", id, next)
			time.Sleep(time.Second)
			runtime.Gosched()
		}
	}
}

func (this *SharedWriter) Commit(id string, lower, upper int64) {
	fmt.Printf("Writer '%s' committing sequence (lower: %d, upper: %d)\n", id, lower, upper)

	originalUpper := upper
	if lower == upper {
		this.committed[upper&this.mask] = int32(upper >> this.shift)
	} else {
		// working down the array rather than up keeps all items in the commit together
		// otherwise the reader(s) could split up the group
		for upper >= lower {
			fmt.Printf("Writer '%s' populating sequence %d (slot %d) with %d\n", id, upper, upper&this.mask, int32(upper>>this.shift))
			this.committed[upper&this.mask] = int32(upper >> this.shift)
			upper--
		}
	}

	fmt.Printf("Writer '%s' committed sequence (lower: %d, upper: %d)\n", id, lower, originalUpper)

}
