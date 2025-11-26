package disruptor

type Disruptor struct {
	Writer
	ReadCloser
}

func NewDisruptor(writer Writer, reader ReadCloser) Disruptor {
	return Disruptor{
		Writer:     writer,
		ReadCloser: reader,
	}
}
