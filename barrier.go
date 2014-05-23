package disruptor

type Barrier interface {
	LoadBarrier(int64) int64
}
