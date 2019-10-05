package disruptor

// Cursors should be a party of the same backing array to keep them as close together as possible:
// https://news.ycombinator.com/item?id=7800825

type Wireup struct {
	capacity int64
	groups   [][]Consumer
	cursors  []*Cursor // backing array keeps cursors (with padding) in contiguous memory
}

func Configure(capacity int64) Wireup {
	return Wireup{capacity: capacity, cursors: []*Cursor{NewCursor()}}
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
	var listeners []ListenCloser
	var upstream Barrier = this.cursors[0]
	written := this.cursors[0]
	cursorIndex := 1 // 0 index is reserved for the writer Cursor

	for groupIndex, group := range this.groups {
		groupReaders, groupBarrier := this.buildReaders(groupIndex, cursorIndex, written, upstream)
		listeners = append(listeners, groupReaders...)
		upstream = groupBarrier
		cursorIndex += len(group)
	}

	sequencer := NewSingleSequencer(written, upstream, this.capacity)
	return Disruptor{sequencer: sequencer, listeners: listeners}
}

func (this Wireup) buildReaders(consumerIndex, cursorIndex int, writeSequence *Cursor, upstream Barrier) ([]ListenCloser, Barrier) {
	var barrierCursors []*Cursor
	var listeners []ListenCloser

	for _, consumer := range this.groups[consumerIndex] {
		readSequence := this.cursors[cursorIndex]
		barrierCursors = append(barrierCursors, readSequence)
		reader := NewReader(readSequence, writeSequence, upstream, consumer)
		listeners = append(listeners, reader)
		cursorIndex++
	}

	if len(barrierCursors) == 0 {
		panic("no barriers")
	}

	if len(this.groups[consumerIndex]) == 1 {
		return listeners, barrierCursors[0]
	} else {
		return listeners, NewCompositeBarrier(barrierCursors)
	}
}
