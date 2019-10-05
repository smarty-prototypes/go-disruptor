package disruptor

import "sync/atomic"

type DefaultReader struct {
	closed   int64
	read     *Sequence // the reader has read up to this sequence
	written  *Sequence // the ring buffer has been written up to this sequence
	upstream Barrier   // the workers just in front of this reader have completed up to this sequence
	consumer Consumer
	waiter   WaitStrategy
}

func NewReader(read, written *Sequence, upstream Barrier, consumer Consumer, waiter WaitStrategy) *DefaultReader {
	return &DefaultReader{
		read:     read,
		written:  written,
		upstream: upstream,
		consumer: consumer,
		waiter:   waiter,
	}
}

func (this *DefaultReader) Listen() {
	current := this.read.Load()
	gateCount, idleCount := 0, 0

	for {
		lower := current + 1
		upper := this.upstream.Load()

		if lower <= upper {
			this.consumer.Consume(lower, upper)
			this.read.Store(upper)
			current = upper
		} else if upper = this.written.Load(); lower <= upper {
			gateCount++
			idleCount = 0
			this.waiter.Gate(gateCount)
		} else if atomic.LoadInt64(&this.closed) > 0 {
			idleCount++
			gateCount = 0
			this.waiter.Idle(idleCount)
		} else {
			break
		}
	}
}

func (this *DefaultReader) Close() error {
	atomic.StoreInt64(&this.closed, 1)
	return nil
}
