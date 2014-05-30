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
	writer := disruptor.NewWriter(written, read, BufferSize)
	reader := disruptor.NewReader(read, written, written, SampleConsumer{})

	started := time.Now()
	reader.Start()
	publish(writer)
	// publish(written, read)
	reader.Stop()
	finished := time.Now()
	fmt.Println(Iterations, finished.Sub(started))
}

func publish(writer *disruptor.Writer) {
	for sequence := disruptor.InitialSequenceValue; sequence <= Iterations; {
		sequence = writer.Reserve()
		ringBuffer[sequence&BufferMask] = sequence
		writer.Commit(sequence)
	}
}

// func publish(written, read *disruptor.Cursor) {
// 	previous := disruptor.InitialSequenceValue
// 	gate := disruptor.InitialSequenceValue

// 	for previous <= Iterations {
// 		next := previous + 1
// 		wrap := next - BufferSize

// 		for wrap > gate {
// 			gate = read.Sequence
// 		}

// 		ringBuffer[next&BufferMask] = next
// 		written.Sequence = next
// 		previous = next
// 	}
// }

type SampleConsumer struct{}

func (this SampleConsumer) Consume(lower, upper int64) {
	for lower <= upper {
		message := ringBuffer[lower&BufferMask]
		if message != lower {
			fmt.Println("Race condition", message, lower)
			panic("Race condition")
		}
		lower++
	}
}
