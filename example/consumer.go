package main

import (
	"time"

	"github.com/smartystreets/go-disruptor"
)

func consume(barrier *disruptor.Barrier, source, sequence *disruptor.Sequence) {
	worker := disruptor.NewWorker(barrier, &ConsumerHandler{time.Now()}, source, sequence)

	for {
		worker.Process()
	}
}
