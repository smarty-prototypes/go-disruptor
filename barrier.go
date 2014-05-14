package main

func (this Barrier) Load(consumer int) int64 {
	minimum := MaxSequenceValue

	for i := 0; i < len(this); i++ {
		cursor := this[i].Load()

		// if len(this) > 1 {
		// 	fmt.Printf("Producer:: Consumer %d Barrier: %d\n", i+1, cursor)
		// } else {
		// 	fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Producer Barrier: %d\n", consumer, cursor)
		// }

		if cursor < minimum {
			minimum = cursor
		}
	}

	return minimum
}

func NewBarrier(upstream ...*Sequence) Barrier {
	this := Barrier{}
	for i := 0; i < len(upstream); i++ {
		this = append(this, upstream[i])
	}
	return this
}

type Barrier []*Sequence
