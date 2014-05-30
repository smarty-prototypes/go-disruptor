package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/smartystreets/go-disruptor"
)

const (
	BufferSize = 1024 * 64
	BufferMask = BufferSize - 1
	Iterations = 1000000 * 100
)

var ringBuffer = [BufferSize]int64{}

func main() {
	runtime.GOMAXPROCS(2)

	written, read := disruptor.NewCursor(), disruptor.NewCursor()

	// reader := disruptor.NewReader(read, written, written, SampleConsumer{})

	started := time.Now()

	// go reader.Start()
	go consume(written, read)
	publish(written, read)
	// reader.Stop()
	// consume(written, read)

	finished := time.Now()
	fmt.Println(Iterations, finished.Sub(started))
}

func publish(written, read *disruptor.Cursor) {
	previous := disruptor.InitialSequenceValue
	gate := disruptor.InitialSequenceValue

	for previous <= Iterations {
		next := previous + 1
		wrap := next - BufferSize

		for wrap > gate {
			gate = read.Sequence
			// if wrap > gate {
			// 	// time.Sleep(time.Nanosecond)
			// 	time.Sleep(time.Microsecond)
			// }
		}

		ringBuffer[next&BufferMask] = next
		written.Sequence = next
		previous = next
	}
}

// func consume(reader *disruptor.Reader) {
func consume(written, read *disruptor.Cursor) {
	sleeps := 0
	consumer := SampleConsumer{}

	previous := int64(-1)
	gate := int64(-1)
	for previous < Iterations {
		current := previous + 1
		gate = written.Sequence

		if current <= gate {

			for current < gate {
				current += consumer.Consume(current, gate)
			}
			// for current <= gate {
			// 	if ringBuffer[current&BufferMask] > 0 {
			// 	}

			// 	current++
			// }

			previous = gate
			read.Sequence = gate
		} else {
			sleeps++
			time.Sleep(time.Microsecond)
		}
	}

	fmt.Println("Consumer sleeps:", sleeps)

	// 	// for sequence, gate := int64(0), int64(0); sequence < Iterations; sequence++ {

	// 	// 	// for gate <= sequence {
	// 	// 	// 	gate = written.Sequence
	// 	// 	// 	if gate <= sequence {
	// 	// 	// 		time.Sleep(time.Microsecond)
	// 	// 	// 	}
	// 	// 	// }

	// 	// 	// if ringBuffer[sequence&BufferMask] > 0 {
	// 	// 	// }

	// 	// 	// read.Sequence = sequence
	// 	// }
}

// // type Consumer interface {
// // 	Consume(lower, upper int64)
// // }

type SampleConsumer struct{}

func (this SampleConsumer) Consume(current, gate int64) int64 {
	if ringBuffer[current&BufferMask] > 0 {
	}

	return 1
}
