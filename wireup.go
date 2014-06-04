package disruptor

// Cursors should be a party of the same backing array to keep them as close together as possible:
// https://news.ycombinator.com/item?id=7800825
type (
	Wireup struct {
		capacity int64
		groups   [][]Consumer
		cursors  []*Cursor // backing array keeps cursors (with padding) in contiguous memory
	}
)

func Configure(capacity int64, consumers ...Consumer) Wireup {
	this := Wireup{
		capacity: capacity,
		groups:   [][]Consumer{},
		cursors:  []*Cursor{NewCursor()},
	}

	return this.WithConsumerGroup(consumers...)
}

func (this Wireup) WithConsumerGroup(consumers ...Consumer) Wireup {
	if len(consumers) == 0 {
		return this
	}

	target := make([]Consumer, len(consumers))
	copy(target, consumers)

	for i := 0; i < len(consumers); i++ {
		this.cursors = append(this.cursors, NewCursor())
	}

	this.groups = append(this.groups, target)
	return this
}

func (this Wireup) Build() Disruptor {
	return NewDisruptor(this)
}

func (this Wireup) BuildShared() SharedDisruptor {
	return NewSharedDisruptor(this)
}
