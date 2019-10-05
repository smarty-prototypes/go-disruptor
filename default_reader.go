package disruptor

type DefaultReader struct {
	closed   *Cursor
	read     *Cursor // the reader has read up to this sequence
	written  *Cursor // the ring buffer has been written up to this sequence
	upstream Barrier // the workers just in front of this reader have completed up to this sequence
	consumer Consumer
	waiter   WaitStrategy
}

func NewReader(read, written *Cursor, upstream Barrier, consumer Consumer, waiter WaitStrategy) *DefaultReader {
	return &DefaultReader{
		closed:   NewCursor(),
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
		} else if this.closed.Load() == InitialCursorSequenceValue {
			idleCount++
			gateCount = 0
			this.waiter.Idle(idleCount)
		} else {
			break
		}
	}
}

func (this *DefaultReader) Close() error {
	this.closed.Store(InitialCursorSequenceValue + 1)
	return nil
}
