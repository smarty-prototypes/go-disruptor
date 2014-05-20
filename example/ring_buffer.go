package main

const (
	RingSize = 1024 * 16
	RingMask = RingSize - 1
)

var ringBuffer [RingSize]int64
