package disruptor

type Disruptor struct {
	Writer
	ListenCloser
}

func NewDisruptor(writer Writer, reader ListenCloser) Disruptor {
	return Disruptor{
		Writer:       writer,
		ListenCloser: reader,
	}
}
