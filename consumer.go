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

	current := uint64(0)
	maxPublished := uint64(0)
	remaining := uint32(0)

	for {
		for current >= maxPublished {
			if !this.started {
				return
			}

			maxPublished = this.producer.AtomicLoad()
			remaining = uint32(maxPublished - current - 1)
			time.Sleep(this.waitStrategy)
		}

		this.handler.Handle(current, remaining)
		this.sequence.Store(current)
		//fmt.Println("Completed sequence:", this.sequence[0])
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
		sequence:     NewSequence(),
		handler:      handler,
		waitStrategy: waitStrategy,
		started:      false,
	}
}
