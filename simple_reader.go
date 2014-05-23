package disruptor

type SimpleReader struct {
	reader   *Reader
	consumer Consumer
}

func NewSimpleReader(reader *Reader, consumer Consumer) *SimpleReader {
	return &SimpleReader{reader: reader, consumer: consumer}
}

func (this *SimpleReader) Receive() (int64, int64) {
	lower, upper := this.reader.Receive()

	if lower <= upper {
		for sequence := lower; sequence <= upper; sequence++ {
			this.consumer.Consume(sequence, upper-sequence)
		}

		this.reader.Commit(lower, upper)
		return lower, upper
	} else {
		return InitialSequenceValue, upper // Idling, Gating
	}
}
