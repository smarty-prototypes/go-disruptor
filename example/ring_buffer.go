package main

const (
	RingSize = 1024 * 1024
	RingMask = RingSize - 1
)

var ringBuffer [RingSize]int64
