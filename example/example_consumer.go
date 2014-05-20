package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/go-disruptor"
)

const Mod = 1000000 * 10 // 1 million * N

func consume(reader *disruptor.Reader) {
	easyReader := disruptor.NewEasyReader(reader, NewSampleConsumer())

	for {
		easyReader.Receive()
	}

	// started := time.Now()
	// consumer := &MyConsumer{started}

	// for {
	// 	sequence, remaining := reader.Receive()
	// 	if remaining >= 0 {
	// 		for remaining >= 0 {
	// 			consumer.Consume(sequence, remaining)
	// 			remaining--
	// 			sequence++
	// 		}

	// 		reader.Commit(sequence)

	// 	} else {
	// 	}
	// }
}

type SampleConsumer struct {
	started time.Time
}

func NewSampleConsumer() disruptor.Consumer {
	return &SampleConsumer{started: time.Now()}
}

func (this *SampleConsumer) Consume(sequence, remaining int64) {
	if sequence%Mod == 0 {
		finished := time.Now()
		fmt.Println(sequence, finished.Sub(this.started))
		this.started = time.Now()
	}

	if sequence != ringBuffer[sequence&RingMask] {
		message := ringBuffer[sequence&RingMask]
		panic(fmt.Sprintf("Sequence: %d, Message %d\n", sequence, message))
	}
}
