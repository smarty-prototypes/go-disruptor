package disruptor

type Consumer interface {
	Consume(lower, upper int64)
}

type Barrier interface {
	Read(int64) int64
}
