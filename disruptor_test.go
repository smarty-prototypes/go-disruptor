package disruptor

import (
	"sync"
	"testing"
)

func TestEndToEnd_SingleWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test")
	}
	coordinator := newEndToEndDisruptor(1, t)
	defer coordinator.Listen()

	go func() {
		defer func() { _ = coordinator.Close() }()
		for i := int64(0); i < endToEndSequences; i++ {
			sequence := coordinator.Reserve(1)
			writeEndToEndEntry(sequence)
			coordinator.Commit(sequence, sequence)
		}
	}()
}
func TestEndToEnd_SharedWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test")
	}
	const writerCount = 2
	coordinator := newEndToEndDisruptor(writerCount, t)
	defer coordinator.Listen()

	go func() {
		var waiter sync.WaitGroup
		waiter.Add(writerCount)
		defer func() { waiter.Wait(); _ = coordinator.Close() }()
		for writerIndex := 0; writerIndex < writerCount; writerIndex++ {
			go func() {
				defer waiter.Done()
				for {
					sequence := coordinator.Reserve(1)
					writeEndToEndEntry(sequence)
					coordinator.Commit(sequence, sequence)
					if sequence >= endToEndSequences-1 {
						return
					}
				}
			}()
		}
	}()
}

func newEndToEndDisruptor(writerCount int, t *testing.T) Disruptor {
	value, err := New(
		Options.BufferCapacity(endToEndBufferSize),
		Options.SingleWriter(writerCount <= 1),
		Options.NewHandlerGroup(newEvenSequenceHandler(), newOddSequenceHandler()),
		Options.NewHandlerGroup(newEvenSequenceHandler(), newOddSequenceHandler()),
		Options.NewHandlerGroup(newEvenSequenceHandler(), newOddSequenceHandler()),
		Options.NewHandlerGroup(newEvenSequenceHandler(), newOddSequenceHandler()),
		Options.NewHandlerGroup(&verificationHandler{t: t}),
	)
	if err != nil {
		panic(err)
	}
	return value
}
func writeEndToEndEntry(sequence int64) {
	entry := &endToEndBuffer[sequence&EndToEndBufferMask]
	entry.Sequence = sequence
	entry.Value1 = sequence + 1
	entry.Value2 = sequence + 2
	entry.Value3 = sequence + 3
	entry.Value4 = sequence + 4
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

type verificationHandler struct{ t *testing.T }

func (this *verificationHandler) Handle(lower, upper int64) {
	for sequence := lower; sequence <= upper; sequence++ {
		assertEndToEndEntry(this.t, int(sequence&EndToEndBufferMask))
	}
}
func assertEndToEndEntry(t *testing.T, index int) {
	t.Helper()
	entry := endToEndBuffer[index]
	sequence := entry.Sequence
	assertEndToEndField(t, index, sequence, "Value1", entry.Value1, sequence+1)
	assertEndToEndField(t, index, sequence, "Value2", entry.Value2, sequence+2)
	assertEndToEndField(t, index, sequence, "Value3", entry.Value3, sequence+3)
	assertEndToEndField(t, index, sequence, "Value4", entry.Value4, sequence+4)
}
func assertEndToEndField(t *testing.T, index int, sequence int64, name string, actual, base int64) {
	t.Helper()
	if expected := base * endToEndMultiplier; actual != expected {
		t.Errorf("entry %d (seq %d): %s=%d, want %d", index, sequence, name, actual, expected)
	}
}

var endToEndBuffer [endToEndBufferSize]endToEndValues

const (
	endToEndBufferSize    = 1024 * 16
	EndToEndBufferMask    = endToEndBufferSize - 1
	endToEndSequences     = endToEndBufferSize * 4
	endToEndHandlerGroups = 4
	endToEndMultiplier    = 1 << endToEndHandlerGroups // 16
)
