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
	Iterations = 1000000 * 100 // 1 million * N
)

// var ringBuffer = [BufferSize]int64{}

func main() {
	runtime.GOMAXPROCS(2)

	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	reader := disruptor.NewReader(read, written, written)
	writer := disruptor.NewWriter(written, read, BufferSize)

	started := time.Now()

	go publish(writer)
	consume(reader)

	finished := time.Now()
	fmt.Println(Iterations, finished.Sub(started))
}

func publish(writer *disruptor.Writer) {
	sequence := disruptor.InitialSequenceValue
	for sequence <= Iterations {
		reservation := writer.Reserve()
		if reservation >= 0 {
			// TODO: publish messages here
			writer.Commit(reservation)
			sequence = reservation
		}
	}
}
func consume(reader *disruptor.Reader) {
	sequence := int64(0)
	for sequence < Iterations {
		received := reader.Receive(sequence)
		if received >= 0 {
			// TODO: handle messages here
			sequence = received
			reader.Commit(sequence)
		} else {
			time.Sleep(time.Microsecond)
		}
	}
}
