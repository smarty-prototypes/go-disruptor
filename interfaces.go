package disruptor

// performance impact of Consume(sequence, remaining int64) (consumed int64)?
type Consumer interface {
	Consume(lower, upper int64)
}

type Barrier interface {
	Read(int64) int64
}
