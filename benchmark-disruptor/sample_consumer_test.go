package benchmarks

import "fmt"

type SampleConsumer struct {
	ringBuffer *[RingBufferSize]int64
}

func (this SampleConsumer) Consume(lower, upper int64) {
	for lower <= upper {
		message := this.ringBuffer[lower&RingBufferMask]
		if message != lower {
			panic(fmt.Sprintf("\nRace condition %d %d\n", lower, message))
		}
		lower++
	}
}
