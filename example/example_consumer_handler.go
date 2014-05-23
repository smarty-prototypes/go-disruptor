package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/go-disruptor"
)

type ExampleConsumerHandler struct {
	started time.Time
}

func NewExampleConsumerHandler() disruptor.Consumer {
	return &ExampleConsumerHandler{started: time.Now()}
}

func (this *ExampleConsumerHandler) Consume(sequence, remaining int64) {
	if sequence%ReportingFrequency == 0 {
		finished := time.Now()
		fmt.Println(sequence, finished.Sub(this.started))
		this.started = time.Now()
	}

	if sequence != ringBuffer[sequence&RingMask] {
		message := ringBuffer[sequence&RingMask]
		alert := fmt.Sprintf("Race Condition::Sequence: %d, Message %d\n", sequence, message)
		fmt.Print(alert)
		panic(alert)
	}
}
