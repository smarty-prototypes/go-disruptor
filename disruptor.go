package disruptor

type Disruptor struct {
	writer  *Writer
	readers []*Reader
}

func NewDisruptor(builder Builder) Disruptor {
	// TODO: Cursors should probably all be created at the same time in wireup to keep them as close together as possible:
	// https://news.ycombinator.com/item?id=7800825
	return Disruptor{}
}

func (this Disruptor) Writer() *Writer {
	return this.writer
}

func (this Disruptor) Start() {
	for _, item := range this.readers {
		item.Start()
	}
}

func (this Disruptor) Stop() {
	for _, item := range this.readers {
		item.Stop()
	}
}
