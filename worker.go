package disruptor

import "time"

type Worker struct {
	ringSequence Sequence
	sequence     Sequence
	callback     func(uint64)
	sleep        time.Duration
}

func (this *Worker) Process() {
	for current, max := uint64(0), uint64(0); ; current++ {
		for current >= max {
			max = this.ringSequence.atomicLoad()
			time.Sleep(this.sleep)
		}

		if max == Uint64MaxValue {
			break
		}

		this.callback(current)
		this.sequence.store(current)
	}
}

func NewWorker(ringSequence Sequence, callback func(uint64), sleep time.Duration) Worker {
	return Worker{
		ringSequence: ringSequence,
		sequence:     NewSequence(),
		callback:     callback,
		sleep:        sleep,
	}
}

const Uint64MaxValue uint64 = 0xFFFFFFFFFFFFFFFF
