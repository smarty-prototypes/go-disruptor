package benchmarks

import "fmt"

type SampleConsumer struct {
}

func (this SampleConsumer) Consume(lower, upper int64) {
	for lower <= upper {
		message := ringBuffer[lower&RingBufferMask]
		if message != lower {
			warning := fmt.Sprintf("\nRace condition--Sequence: %d, Message: %d\n", lower, message)
			fmt.Printf(warning)
			panic(warning)
		}
		lower++
	}
}
