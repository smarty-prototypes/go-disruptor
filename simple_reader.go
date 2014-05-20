package disruptor

type SimpleReader struct {
	reader   *Reader
	consumer Consumer
}

func NewSimpleReader(reader *Reader, consumer Consumer) *SimpleReader {
	return &SimpleReader{reader: reader, consumer: consumer}
}

func (this *SimpleReader) Receive() int64 {
	sequence, remaining := this.reader.Receive()

	if remaining >= 0 {
		for remaining >= 0 {
			this.consumer.Consume(sequence, remaining)
			remaining--
			sequence++
		}

		this.reader.Commit(sequence - 1)
		return sequence
	} else {
		return remaining // Idling, Gating
	}
}
