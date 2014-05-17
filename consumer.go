package disruptor

type Consumer interface {
	Consume(sequence, remaining int64)
}
