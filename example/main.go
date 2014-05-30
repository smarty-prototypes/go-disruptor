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

	go consume(written, read, SampleConsumer{})
	publish(written, read)

	finished := time.Now()
	fmt.Println(Iterations, finished.Sub(started))

	time.Sleep(time.Millisecond * 10)
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

func consume(written, read *disruptor.Cursor, consumer disruptor.Consumer) {
	previous := int64(-1)
	gate := int64(-1)
	for previous < Iterations {
		current := previous + 1
		gate = written.Sequence

		if current <= gate {
			consumer.Consume(current, gate)
			previous = gate
			read.Sequence = gate
		} else {
			// TODO: wait strategy
		}

		time.Sleep(time.Nanosecond)
	}
}

type SampleConsumer struct{}

func (this SampleConsumer) Consume(lower, upper int64) {
	for lower <= upper {
		if ringBuffer[lower&BufferMask] > 0 {
		}
		lower++
	}
}
