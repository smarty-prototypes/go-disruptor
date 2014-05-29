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
	started := time.Now()

	go publish(written, read)
	consume(written, read)

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
			// 	time.Sleep(time.Microsecond)
			// }
		}

		ringBuffer[next&BufferMask] = next
		written.Sequence = next
		previous = next
	}
}
func consume(written, read *disruptor.Cursor) {
	for sequence, gate := int64(0), int64(0); sequence < Iterations; sequence++ {
		for gate <= sequence {
			gate = written.Sequence
			if gate <= sequence {
				time.Sleep(time.Microsecond)
			}
		}

		if ringBuffer[sequence&BufferMask] > 0 {
		}

		read.Sequence = sequence
	}
}
