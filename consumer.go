package disruptor

type Consumer struct {
	sequence Sequence
	handler  MessageHandler
	started  bool
}

func (this *Consumer) Start() {
}

func (this *Consumer) Stop() {
	this.started = false
}

func NewConsumer(handler MessageHandler) Consumer {
	return Consumer{handler: handler}
}
