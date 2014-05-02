package disruptor

type MessageHandler interface {
	Handle(sequence uint64, remaining uint32)
}
