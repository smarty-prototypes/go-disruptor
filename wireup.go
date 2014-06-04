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
	allReaders := []*Reader{}
	written := this.cursors[0]
	barriers := []Barrier{written}
	cursorIndex := 1 // 0 index is reserved for the writer Cursor

	for groupIndex, group := range this.groups {
		upstream := barriers[groupIndex]
		groupReaders, barrier := this.buildReaders(groupIndex, cursorIndex, written, upstream)
		for _, item := range groupReaders {
			allReaders = append(allReaders, item)
		}
		barriers = append(barriers, barrier)
		cursorIndex += len(group)
	}

	writerBarrier := barriers[len(barriers)-1]
	writer := NewWriter(written, writerBarrier, this.capacity)
	return Disruptor{writer: writer, readers: allReaders}
}
func (this Wireup) buildReaders(consumerIndex, cursorIndex int, written *Cursor, upstream Barrier) ([]*Reader, Barrier) {
	barrierCursors := []*Cursor{}
	readers := []*Reader{}

	for _, consumer := range this.groups[consumerIndex] {
		cursor := this.cursors[cursorIndex]
		barrierCursors = append(barrierCursors, cursor)
		reader := NewReader(cursor, written, upstream, consumer)
		readers = append(readers, reader)
		cursorIndex++
	}

	if len(this.groups[consumerIndex]) == 1 {
		return readers, barrierCursors[0]
	} else {
		return readers, NewCompositeBarrier(barrierCursors...)
	}
}
