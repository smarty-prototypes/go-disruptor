package disruptor

type Disruptor struct {
	writer  Writer
	readers []*Reader
}

func (this Disruptor) Listen() {
	for _, item := range this.readers {
		go item.Listen()
	}
}

func (this Disruptor) Close() error {
	for _, item := range this.readers {
		_ = item.Close()
	}

	return nil
}

func (this Disruptor) Writer() Writer { return this.writer }
