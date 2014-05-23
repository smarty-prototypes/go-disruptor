package disruptor

type SimpleReader struct {
	reader   *Reader
	callback Consumer
}

func NewSimpleReader(reader *Reader, callback Consumer) *SimpleReader {
	return &SimpleReader{reader: reader, callback: callback}
}

func (this *SimpleReader) Receive(lower int64) int64 {
	upper := this.reader.Receive(lower)

	if lower <= upper {
		for sequence := lower; sequence <= upper; sequence++ {
			this.callback.Consume(sequence, upper-sequence)
		}

		this.reader.Commit(upper)
	}

	return upper
}
