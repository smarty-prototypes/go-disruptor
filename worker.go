package disruptor

func (this *Worker) Process() int64 {
	next := this.sequence.Load() + 1
	ready := this.barrier.Load()

	if next <= ready {
		for next <= ready {
			this.handler.Consume(next, ready-next)
			next++
		}

		next--
		this.sequence.Store(next)
		return next
	} else if next <= this.source.Load() {
		return Gating
	} else {
		return Idle
	}
}

func NewWorker(barrier Barrier, handler Consumer, source, sequence *Sequence) *Worker {
	return &Worker{
		barrier:  barrier,
		handler:  handler,
		source:   source,
		sequence: sequence,
	}
}

type Worker struct {
	barrier  Barrier
	handler  Consumer
	source   *Sequence
	sequence *Sequence
}

type Consumer interface {
	Consume(sequence, remaining int64)
}

const (
	Gating int64 = -2
	Idle         = -3
)
