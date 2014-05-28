package disruptor

type Consumer interface {
	Consume(lower, upper int64)
}
