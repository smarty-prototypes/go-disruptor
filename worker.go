package disruptor

import "time"

type Worker struct {
	cursor   Sequence
	sequence Sequence
	callback func(uint64)
	sleep    time.Duration
}

func (this *Worker) Process() {
	current, max := uint64(0), uint64(0)
	for {
		for current >= max {
			max = this.cursor.atomicLoad()
			time.Sleep(this.sleep)
		}

		if max == Uint64MaxValue {
			break
		}

		this.callback(current)
		this.sequence.store(current)
		current++
	}
}

func NewWorker(cursor Sequence, callback func(uint64), sleep time.Duration) Worker {
	return Worker{
		cursor:   cursor,
		sequence: NewSequence(),
		callback: callback,
		sleep:    sleep,
	}
}

const Uint64MaxValue uint64 = 0xFFFFFFFFFFFFFFFF
