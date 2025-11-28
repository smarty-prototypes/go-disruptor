package disruptor

import "sync"

type compositeListener []ListenCloser

func (this compositeListener) Listen() {
	waiter := &sync.WaitGroup{}
	waiter.Add(len(this))
	defer waiter.Wait()

	for _, item := range this {
		go func(listener ListenCloser) {
			listener.Listen()
			waiter.Done()
		}(item)
	}
}

func (this compositeListener) Close() error {
	for _, item := range this {
		_ = item.Close()
	}

	return nil
}
