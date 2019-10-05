package disruptor

import "time"

type DefaultWaitStrategy struct{}

func NewWaitStrategy() DefaultWaitStrategy      { return DefaultWaitStrategy{} }
func (this DefaultWaitStrategy) Gate(count int) { time.Sleep(time.Microsecond) }
func (this DefaultWaitStrategy) Idle(count int) { time.Sleep(time.Millisecond) }
