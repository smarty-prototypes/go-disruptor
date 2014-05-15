package main

import "fmt"

type TestHandler struct{}

func (this TestHandler) Consume(sequence, remaining int64) {
	message := ringBuffer[sequence&RingMask]

	if message != sequence {
		text := fmt.Sprintf("[Consumer] ERROR Sequence: %d, Message: %d\n", sequence, message)
		fmt.Printf(text)
		panic(text)
	}

	if sequence%Mod == 0 && sequence > 0 {
		// fmt.Printf("[Consumer] Sequence: %d, Message: %d\n", sequence, message)
	}
}
