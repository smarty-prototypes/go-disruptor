package disruptor

type Consumer interface {
	Consume(int64, int64)
}
