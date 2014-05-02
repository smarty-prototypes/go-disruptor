package disruptor

type MessageHandler interface {
	Handle(uint64 sequence, uint32 remaining)
}
