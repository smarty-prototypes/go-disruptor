package main

type ConsumerHandler struct{}

func (this ConsumerHandler) Consume(sequence, remaining int64) {
	message := ringBuffer[sequence&RingMask]
	if message != sequence {
		panic("Race condition")
	}
}
