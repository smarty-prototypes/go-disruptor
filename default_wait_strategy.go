package disruptor

import "time"

type defaultWaitStrategy struct{}

func (this defaultWaitStrategy) Gate(int64) { time.Sleep(time.Nanosecond) }
func (this defaultWaitStrategy) Idle(int64) { time.Sleep(time.Millisecond) }
