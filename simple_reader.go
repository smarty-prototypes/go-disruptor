package disruptor

type SimpleReader struct {
	reader   *Reader
	callback Consumer
}

func NewSimpleReader(reader *Reader, callback Consumer) *SimpleReader {
	return &SimpleReader{reader: reader, callback: callback}
}

func (this *SimpleReader) Receive() (int64, int64) {
	lower, upper := this.reader.Receive()

	if lower <= upper {
		for sequence := lower; sequence <= upper; sequence++ {
			this.callback.Consume(sequence, upper-sequence)
		}

		this.reader.Commit(lower, upper)
		return lower, upper
	} else {
		return InitialSequenceValue, upper // Idling, Gating
	}
}
