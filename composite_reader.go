package disruptor

import "sync"

type compositeReader []ReadCloser

func (this compositeReader) Read() {
	var waiter sync.WaitGroup
	waiter.Add(len(this))

	for _, item := range this {
		go func(reader ReadCloser) {
			reader.Read()
			waiter.Done()
		}(item)
	}

	waiter.Wait()
}

func (this compositeReader) Close() error {
	for _, item := range this {
		_ = item.Close()
	}

	return nil
}
