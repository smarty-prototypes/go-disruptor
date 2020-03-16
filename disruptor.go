package disruptor

type Disruptor struct {
	Writer
	Reader
}

func NewDisruptor(writer Writer, reader Reader) Disruptor {
	return Disruptor{
		Writer: writer,
		Reader: reader,
	}
}
