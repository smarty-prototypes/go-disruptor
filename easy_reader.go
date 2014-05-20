package disruptor

type EasyReader struct {
	reader   *Reader
	consumer Consumer
}

func NewEasyReader(reader *Reader, consumer Consumer) *EasyReader {
	return &EasyReader{reader: reader, consumer: consumer}
}

func (this *EasyReader) Receive() int64 {
	sequence, remaining := this.reader.Receive()

	if remaining >= 0 {
		for remaining >= 0 {
			this.consumer.Consume(sequence, remaining)
			remaining--
			sequence++
		}

		this.reader.Commit(sequence)
		return sequence
	} else {
		return remaining // Idling, Gating
	}
}
