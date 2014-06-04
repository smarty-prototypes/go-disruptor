package disruptor

type (
	Wireup struct {
		capacity int64
		groups   [][]Consumer
	}
)

func Configure(capacity int64, consumers ...Consumer) Wireup {
	this := Wireup{
		capacity: capacity,
		groups:   [][]Consumer{},
	}

	return this.WithConsumerGroup(consumers...)
}

func (this Wireup) WithConsumerGroup(consumers ...Consumer) Wireup {
	if len(consumers) == 0 {
		return this
	}

	target := make([]Consumer, len(consumers))
	copy(target, consumers)

	this.groups = append(this.groups, target)
	return this
}

func (this Wireup) Build() Disruptor {
	return NewDisruptor(this)
}

func (this Wireup) BuildShared() SharedDisruptor {
	return NewSharedDisruptor(this)
}
