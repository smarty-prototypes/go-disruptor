package disruptor

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestEndToEnd_SingleWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test")
	}
	coordinator := newEndToEndDisruptor(1)

	go func() {
		defer func() { _ = coordinator.Close() }()
		for i := int64(0); i < endToEndSequences; i++ {
			sequence := coordinator.Reserve(1)
			writeEndToEndEntry(sequence)
			coordinator.Commit(sequence, sequence)
		}
	}()

	coordinator.Listen()
	verifyEndToEndBuffer(t)
}
func TestEndToEnd_SharedWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test")
	}
	coordinator := newEndToEndDisruptor(endToEndWriterCount)

	go func() {
		var waiter sync.WaitGroup
		waiter.Add(endToEndWriterCount)
		defer func() { waiter.Wait(); _ = coordinator.Close() }()

		var remaining atomic.Int64
		remaining.Store(endToEndSequences)
		for writerIndex := 0; writerIndex < endToEndWriterCount; writerIndex++ {
			go func() {
				defer waiter.Done()
				for remaining.Add(-1) >= 0 {
					sequence := coordinator.Reserve(1)
					writeEndToEndEntry(sequence)
					coordinator.Commit(sequence, sequence)
				}
			}()
		}
	}()

	coordinator.Listen()
	verifyEndToEndBuffer(t)
}

func newEndToEndDisruptor(writerCount int) Disruptor {
	endToEndBuffer = [endToEndBufferSize]endToEndValues{}
	value, err := New(
		Options.BufferCapacity(endToEndBufferSize),
		Options.SingleWriter(writerCount <= 1),
		Options.NewHandlerGroup(newEvenSequenceHandler(), newOddSequenceHandler()),
		Options.NewHandlerGroup(newEvenSequenceHandler(), newOddSequenceHandler()),
	)
	if err != nil {
		panic(err)
	}
	return value
}
func writeEndToEndEntry(sequence int64) {
	entry := &endToEndBuffer[sequence&EndToEndBufferMask]
	entry.Sequence = sequence
	entry.Value1 += sequence + 1
	entry.Value2 += sequence + 2
	entry.Value3 += sequence + 3
	entry.Value4 += sequence + 4
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type endToEndHandler bool
type endToEndValues struct {
	Sequence int64    // 8B
	Value1   int64    // 8B
	Value2   int64    // 8B
	Value3   int64    // 8B
	Value4   int64    // 8B
	_        [3]int64 // 24B padding to 64B cache line
}

func newEvenSequenceHandler() Handler { return endToEndHandler(true) }
func newOddSequenceHandler() Handler  { return endToEndHandler(false) }
func (this endToEndHandler) Handle(lower, upper int64) {
	for sequence := lower; sequence <= upper; sequence++ {
		if (sequence&1 == 0) != this {
			continue
		}

		entry := &endToEndBuffer[sequence&EndToEndBufferMask]
		entry.Value1 *= 2
		entry.Value2 *= 2
		entry.Value3 *= 2
		entry.Value4 *= 2
	}
}

func verifyEndToEndBuffer(t *testing.T) {
	t.Helper()
	for i := range endToEndBuffer {
		assertEndToEndEntry(t, i)
	}
}
func assertEndToEndEntry(t *testing.T, index int) {
	t.Helper()
	entry := endToEndBuffer[index]
	assertEndToEndField(t, index, "Value1", entry.Value1, expectedEndToEndValue(index, 1))
	assertEndToEndField(t, index, "Value2", entry.Value2, expectedEndToEndValue(index, 2))
	assertEndToEndField(t, index, "Value3", entry.Value3, expectedEndToEndValue(index, 3))
	assertEndToEndField(t, index, "Value4", entry.Value4, expectedEndToEndValue(index, 4))
}
func expectedEndToEndValue(index int, offset int64) int64 {
	var value int64
	for wrap := int64(0); wrap < endToEndWraps; wrap++ {
		sequence := int64(index) + wrap*endToEndBufferSize
		value = (value + sequence + offset) * endToEndMultiplier
	}
	return value
}
func assertEndToEndField(t *testing.T, index int, name string, actual, expected int64) {
	t.Helper()
	if actual != expected {
		t.Errorf("entry %d: %s=%d, want %d", index, name, actual, expected)
	}
}

var endToEndBuffer [endToEndBufferSize]endToEndValues

const (
	endToEndBufferSize      = 1024 * 128
	EndToEndBufferMask      = endToEndBufferSize - 1
	endToEndSequences       = endToEndBufferSize * 4
	endToEndWraps           = endToEndSequences / endToEndBufferSize
	endToEndHandlerGroups   = 2
	endToEndMultiplier  = 1 << endToEndHandlerGroups
	endToEndWriterCount = 2
)
