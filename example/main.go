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
	Iterations = 1000000 * 1000
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

		previous = next

		ringBuffer[next&BufferMask] = next

		written.Sequence = next
	}

	// next := this.previous + 1
	// wrap := next - this.capacity

	// for wrap > this.gate {
	// 	this.gate = this.upstream.Read(next)
	// }

	// this.previous = next
	// return next

	// // gates := 0
	// sequence := disruptor.InitialSequenceValue
	// for sequence <= Iterations {
	// 	sequence = writer.Reserve()
	// 	written.Sequence = sequence
	// 	// writer.Commit(sequence)
	// 	// reservation := writer.Reserve()
	// 	// if reservation >= 0 {
	// 	// 	// TODO: publish messages here
	// 	// 	writer.Commit(reservation)
	// 	// 	sequence = reservation
	// 	// } else {
	// 	// 	gates++
	// 	// }
	// }

	// fmt.Println("Write gates", gates)
}
func consume(written, read *disruptor.Cursor) {
	// gates := 0
	sequence := int64(0)
	for sequence < Iterations {
		maximum := written.Sequence
		for maximum <= sequence {
			// gates++
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

	// fmt.Println("Total Read gates", gates)
}
