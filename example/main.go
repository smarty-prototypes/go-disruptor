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
		}

		ringBuffer[next&BufferMask] = next
		written.Sequence = next
		previous = next
	}
}
func consume(written, read *disruptor.Cursor) {
	sequence := int64(0)
	for sequence < Iterations {
		maximum := written.Sequence
		for maximum <= sequence {
			maximum = written.Sequence
			time.Sleep(time.Microsecond)
		}

		for sequence < maximum {
			if ringBuffer[sequence&BufferMask] > 0 {
			}
			sequence++
		}

		read.Sequence = maximum
		sequence = maximum + 1
	}
}
