package benchmarks

import "fmt"

type SampleConsumer struct {
	ringBuffer *[RingBufferSize]int64
}

func (this SampleConsumer) Consume(lower, upper int64) {
	for lower <= upper {
		message := this.ringBuffer[lower&RingBufferMask]
		if message != lower {
			warning := fmt.Sprintf("\nRace condition--Sequence: %d, Message: %d\n", lower, message)
			fmt.Printf(warning)
			panic(warning)
		}
		lower++
	}
}
