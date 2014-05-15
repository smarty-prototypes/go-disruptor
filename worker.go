package main

func (this Worker) Process(consumer int) int64 {
	next := this.sequence.Load() + 1
	ready := this.barrier.Load(consumer)

	if next <= ready {
		// fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: %d work items found (Next %d, Ready %d).\n", consumer, ready-next+1, ready, next)

		// if ready > 1000 {
		// 	fmt.Println("\t\t\t\t\t\t\t\t\tMalformed Next!", next, ready)
		// }

		for next <= ready {
			// fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Consuming sequence %d\n", consumer, next)
			this.handler.Consume(next, ready-next)
			next++
		}

		// fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Completed through sequence %d\n", consumer, next-1)
		next--
		this.sequence.Store(next)
		return next
	} else if next <= this.source.Load() {
		return Gating
	} else {
		return Idle
	}
}

func NewWorker(barrier Barrier, handler Consumer, source, sequence *Sequence) Worker {
	return Worker{
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
