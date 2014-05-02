package disruptor

import "time"

type Consumer struct {
	producer     Sequence
	sequence     Sequence
	handler      MessageHandler
	waitStrategy time.Duration
	started      bool
}

func (this *Consumer) Start() {
	this.started = true
	wait := this.waitStrategy
	producer := this.producer
	sequence := this.sequence
	handler := this.handler

	current := uint64(0)
	maxPublished := uint64(0)
	remaining := uint32(0)

	for {
		for current >= maxPublished {
			if !this.started {
				return
			}

			maxPublished = producer.AtomicLoad()
			remaining = uint32(maxPublished - current)
			time.Sleep(wait)
		}

		handler.Handle(current, remaining)
		sequence.Store(current)
		current++
		remaining--
	}
}

func (this *Consumer) Stop() {
	this.started = false
}

func NewConsumer(producer Sequence, handler MessageHandler, waitStrategy time.Duration) Consumer {
	return Consumer{
		producer:     producer,
		handler:      handler,
		waitStrategy: waitStrategy,
	}
}
