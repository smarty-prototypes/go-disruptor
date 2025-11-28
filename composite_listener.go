package disruptor

import "sync"

type compositeListener []ListenCloser

func newCompositeListener(listeners []ListenCloser) ListenCloser {
	if len(listeners) == 1 {
		return listeners[0]
	} else {
		return compositeListener(listeners)
	}
}

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
