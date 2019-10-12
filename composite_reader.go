package disruptor

import "sync"

type compositeReader []Reader

func (this compositeReader) Read() {
	var waiter sync.WaitGroup
	waiter.Add(len(this))

	for _, item := range this {
		go func(reader Reader) {
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
