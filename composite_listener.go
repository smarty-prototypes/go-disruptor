package disruptor

import "sync"

// compositeListener is a ListenCloser that manages multiple listeners concurrently. Listen starts each child
// listener on its own goroutine and waits for all of them to complete. Close propagates to all child listeners.
// The constructor optimizes for the single-listener case by returning it unwrapped.
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
			defer waiter.Done()
			listener.Listen()
		}(item)
	}
}

func (this compositeListener) Close() error {
	for _, item := range this {
		_ = item.Close()
	}

	return nil
}
