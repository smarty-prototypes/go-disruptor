package disruptor

import "runtime"

type SpinSequencer struct{ Sequencer }

func NewSpinSequencer(inner Sequencer) *SpinSequencer { return &SpinSequencer{Sequencer: inner} }

func (this SpinSequencer) Reserve(count int64) int64 {
	for spinCount := int64(0); ; spinCount++ {
		if sequence := this.Sequencer.Reserve(count); sequence > defaultSequenceValue {
			return sequence
		} else if spinCount > 0 && spinCount&spinMask == 0 {
			runtime.Gosched() // http://bit.ly/1xiDINZ
		}
	}
}

const spinMask = 1024*16 - 1 // TODO: experiment with different values (always a power of 2 and then subtract 1)
