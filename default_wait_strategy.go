package disruptor

import "time"

type DefaultWaitStrategy struct{}

func NewWaitStrategy() DefaultWaitStrategy        { return DefaultWaitStrategy{} }
func (this DefaultWaitStrategy) Gate(count int64) { time.Sleep(time.Nanosecond) }
func (this DefaultWaitStrategy) Idle(count int64) { time.Sleep(time.Microsecond * 50) }
