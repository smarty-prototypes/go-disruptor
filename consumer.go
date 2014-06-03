package disruptor

// TODO: investigate performance impact of Consume(sequence, remaining int64) (consumed int64)
type Consumer interface {
	Consume(lower, upper int64)
}
