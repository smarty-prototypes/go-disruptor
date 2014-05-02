package disruptor

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkDisruptor(b *testing.B) {
	//consumer, producer := NewSequence(), NewSequence()
	//ring := NewRingBuffer(BufferSize)
	iterations := uint64(b.N)
	iterations = 64

	producerSequence := NewSequence()
	handler := MyHandler{}
	consumer := NewConsumer(producerSequence, handler, WaitStrategy)
	consumerSequence := consumer.sequence

	go func() {
		for current, maxAvailable := uint64(0), uint64(0); current < iterations; {
			for current >= maxAvailable {
				maxAvailable = consumerSequence.AtomicLoad() + BufferSize
				//fmt.Println("Max available", maxAvailable, iterations)
				time.Sleep(WaitStrategy)
			}

			//ring[current&BufferMask] = current
			current++
			producerSequence.Store(current)
		}
	}()

	consumer.Start()

	// for current, maxPublished := uint64(0), uint64(0); current < iterations; current++ {
	// 	for current >= maxPublished {
	// 		maxPublished = producer.AtomicLoad()
	// 		time.Sleep(WaitStrategy)
	// 	}

	// 	message := ring[current&BufferMask]
	// 	if message != current {
	// 		panic("Out of sequence")
	// 	}
	// 	consumer.Store(current)
	// }
}

type MyHandler struct{}

func (this MyHandler) Handle(sequence uint64, remaining uint32) {
	fmt.Println("Current Sequence:", sequence, remaining)
	//time.Sleep(time.Millisecond * 250)
}

const BufferSize = 1024 * 128
const BufferMask = BufferSize - 1
const WaitStrategy = time.Nanosecond
