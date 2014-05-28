package disruptor

type BatchReader struct {
	started  bool
	read     *Cursor
	written  *Cursor
	upstream Barrier
	consumer Consumer
	waiter   Waiter
}

func NewBatchReader(read, written *Cursor, upstream Barrier, consumer Consumer, waiter Waiter) *BatchReader {
	return &BatchReader{
		started:  false,
		read:     read,
		written:  written,
		upstream: upstream,
		consumer: consumer,
		waiter:   waiter,
	}
}

func (this *BatchReader) Start() {
	this.started = true
	go this.receive()
}
func (this *BatchReader) Stop() {
	this.started = false
}

func (this *BatchReader) receive() {
	sequence := this.read.Load()

	for {
		lower := sequence + 1
		upper := this.upstream.LoadBarrier(lower)

		if lower <= upper {
			this.consumer.Consume(lower, upper)
			sequence = upper
			this.read.Store(sequence)
		} else if gate := this.written.Load(); lower <= gate {
			// time.Sleep(time.Millisecond) // TODO: use another method from the wait strategy?
		} else if this.started {
			// this.waiter.Wait()
		} else {
			break
		}
	}
}
