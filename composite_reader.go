package disruptor

import "sync"

type compositeReader []ListenCloser

func (this compositeReader) Listen() {
	var waiter sync.WaitGroup
	waiter.Add(len(this))

	for _, item := range this {
		go func(reader ListenCloser) {
			reader.Listen()
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
