package disruptor

type SharedDisruptor struct {
	writer  *Writer
	readers []*Reader
}

func NewSharedDisruptor(builder Builder) SharedDisruptor {
	return SharedDisruptor{}
}

func (this SharedDisruptor) Writer() *Writer {
	return this.writer
}

func (this SharedDisruptor) Start() {
	for _, item := range this.readers {
		item.Start()
	}
}

func (this SharedDisruptor) Stop() {
	for _, item := range this.readers {
		item.Stop()
	}
}
