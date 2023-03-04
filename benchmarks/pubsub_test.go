package benchmarks

import (
	"context"
	"sync"
	"testing"

	"github.com/smartystreets-prototypes/go-disruptor"
	"github.com/stretchr/testify/assert"
)

type RingBuffer[A any] struct {
	buf     []A
	bufMask int64
	rCount  int // reserveCount
	dis     disruptor.Disruptor
}

func NewRingBuffer[A any](ringBufferSize int64, consumeGroup [][]func(a A)) *RingBuffer[A] {
	buf := make([]A, ringBufferSize)
	bufMask := ringBufferSize - 1

	rb := &RingBuffer[A]{
		buf:     buf,
		bufMask: bufMask,
		rCount:  4,
	}

	options := make([]disruptor.Option, 0, len(consumeGroup)+1)
	options = append(options, disruptor.WithCapacity(ringBufferSize))

	for _, group := range consumeGroup {
		cs := make([]disruptor.Consumer, 0, len(group))

		for _, f := range group {
			cs = append(cs, &RingBufferConsumer[A]{
				Buf:     buf,
				BufMask: bufMask,
				Func:    f,
			})
		}

		options = append(options, disruptor.WithConsumerGroup(cs...))
	}

	dis := disruptor.New(options...)

	rb.dis = dis

	return rb
}

func (r *RingBuffer[A]) Start(ctx context.Context) error {
	r.dis.Read()
	return nil
}

func (r *RingBuffer[A]) Stop(ctx context.Context) error {
	return r.dis.Close()
}

func (r *RingBuffer[A]) Send(ctx context.Context, a ...A) error {
	//fmt.Println(a)
	for i := 0; i < len(a); {
		reserveCount := r.rCount
		if len(a)-i < r.rCount {
			reserveCount = len(a) - i
		}

		upper := r.dis.Reserve(int64(reserveCount))
		lower := upper - (int64(reserveCount) - 1)
		for j := lower; j <= upper; j++ {
			//fmt.Println("sequence", sequence, "j", j, j&r.bufMask)
			r.buf[j&r.bufMask] = a[i]

			i++
		}
		r.dis.Commit(lower, upper)
	}

	return nil
}

type RingBufferConsumer[A any] struct {
	Buf     []A
	BufMask int64
	Func    func(a A)
}

func (c *RingBufferConsumer[A]) Consume(lower, upper int64) {
	var message A
	for lower <= upper {
		message = c.Buf[lower&c.BufMask]

		c.Func(message)

		lower++
	}
}

func TestRingBuffer(t *testing.T) {
	ctx := context.Background()

	g1c1Sum := 0
	g2c1Sum := 0
	g2c2Sum := 0

	rb := NewRingBuffer[int](16, [][]func(a int){
		{
			func(a int) {
				g1c1Sum += a
				t.Log("group-1", "consumer-1", a)
			},
		},
		{
			func(a int) {
				g2c1Sum += a

				t.Log("group-2", "consumer-1", a)
			},
			func(a int) {
				g2c2Sum += a

				t.Log("group-2", "consumer-2", a)
			},
		},
	})

	count := 10
	// 3 producer
	wg := sync.WaitGroup{}
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(idx int) {
			start := idx * count

			for start < (idx+1)*count {

				errLocal := rb.Send(ctx, start, start+1)
				if errLocal != nil {
					t.Error(errLocal)
				}

				start += 2
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		err := rb.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()

	err := rb.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 735, g1c1Sum)
	assert.Equal(t, 735, g2c1Sum)
	assert.Equal(t, 735, g2c2Sum)
}

func BenchmarkRB(b *testing.B) {
	benchmarkRingBuffer(b)
}

func benchmarkRingBuffer(b *testing.B) {
	ctx := context.Background()
	iterations := int64(b.N)
	b.Log(iterations)

	rb := NewRingBuffer[int64](RingBufferSize, [][]func(a int64){
		{
			func(a int64) {
			},
		},
		{
			func(a int64) {
			},
			func(a int64) {
			},
		},
	})

	go func() {
		b.ReportAllocs()
		b.ResetTimer()

		size := 100
		var sequence int64 = 1
		for ; sequence < iterations; sequence++ {
			seqList := make([]int64, 0, size)
			for count := 0; count < size && sequence < iterations; count++ {
				seqList = append(seqList, sequence)
				sequence++
			}
			err := rb.Send(ctx, seqList...)
			if err != nil {
				b.Error(err)
			}
		}

		err := rb.Stop(ctx)
		if err != nil {
			b.Error(err)
		}
	}()

	err := rb.Start(ctx)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkChan(b *testing.B) {
	benchmarkChan(b)
}

func benchmarkChan(b *testing.B) {
	iterations := int64(b.N)
	b.Log(iterations)

	ch := make(chan int64, RingBufferSize)

	go func() {
		b.ReportAllocs()
		b.ResetTimer()

		var sequence int64 = 1
		for ; sequence < iterations; sequence++ {
			ch <- sequence
		}

		close(ch)
	}()

	for i := range ch {
		_ = i
	}
}

/**
TODO
1. 多生产者并发写、有线程安全问题。 需要加上mutex，或者atomic操作。
2. consumerGroup 设计上应该是想，一个group的消费者共享cursor。 但是，实际跑起来，一个group中的消费者，都消费了全量数据。这一点看起来，与目标不符。
*/
