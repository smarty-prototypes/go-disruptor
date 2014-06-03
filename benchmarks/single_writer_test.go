package benchmarks

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/smartystreets/go-disruptor"
)

const (
	singleWriterRingBufferSize = 1024 * 64
	singleWriterRingBufferMask = singleWriterRingBufferSize - 1
	reserveOne                 = 1
	reserveMany                = 16
	reserveManyDelta           = reserveMany - 1
)

func BenchmarkDisruptorWriterReserveSingle(b *testing.B) {
	ringBuffer := [singleWriterRingBufferSize]int64{}
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	reader := disruptor.NewReader(read, written, written, singleWriterConsumer{&ringBuffer})
	writer := disruptor.NewWriter(written, read, singleWriterRingBufferSize)
	reader.Start()

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	sequence := disruptor.InitialSequenceValue
	for sequence < iterations {
		sequence = writer.Reserve(reserveOne)
		ringBuffer[sequence&singleWriterRingBufferMask] = sequence
		writer.Commit(sequence, sequence)
	}

	reader.Stop()
}
func BenchmarkDisruptorWriterReserveMultiple(b *testing.B) {
	ringBuffer := [singleWriterRingBufferSize]int64{}
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	reader := disruptor.NewReader(read, written, written, singleWriterConsumer{&ringBuffer})
	writer := disruptor.NewWriter(written, read, singleWriterRingBufferSize)
	reader.Start()

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	sequence := disruptor.InitialSequenceValue
	for sequence < iterations {
		sequence = writer.Reserve(reserveMany)

		for i := sequence - reserveManyDelta; i <= sequence; i++ {
			ringBuffer[i&singleWriterRingBufferMask] = i
		}

		writer.Commit(sequence, sequence)
	}

	reader.Stop()
}

func benchmarkSingleWriter(b *testing.B, maxClaim int64) {

}

type singleWriterConsumer struct {
	ringBuffer *[singleWriterRingBufferSize]int64
}

func (this singleWriterConsumer) Consume(lower, upper int64) {
	for lower <= upper {
		message := this.ringBuffer[lower&singleWriterRingBufferMask]
		if message != lower {
			panic(fmt.Sprintf("\nRace condition %d %d\n", lower, message))
		}
		lower++
	}
}

func init() {
	runtime.GOMAXPROCS(2)
}
