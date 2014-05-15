package main

const RingSize = 1024 * 128
const RingMask = RingSize - 1

var ringBuffer [RingSize]int64
