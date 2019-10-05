package main

import (
	"log"

	"github.com/smartystreets-prototypes/go-disruptor"
)

func main() {
	wireup, err := disruptor.New(
		disruptor.WithCapacity(32),
		disruptor.WithConsumerGroup(MyConsumer{}))

	if err != nil {
		panic(err)
	}

	sequencer, listener := wireup.Build()

	go func() {
		sequence := sequencer.Reserve(1)
		sequencer.Commit(sequence, sequence)

		sequence = sequencer.Reserve(1)
		sequencer.Commit(sequence, sequence)

		sequence = sequencer.Reserve(15)
		sequencer.Commit(sequence, sequence)

		_ = listener.Close()
	}()

	listener.Listen()
}

type MyConsumer struct{}

func (this MyConsumer) Consume(lower, upper int64) {
	log.Println(lower, upper)
}
