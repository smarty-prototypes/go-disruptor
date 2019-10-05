package disruptor

import "sync"

type CompositeListenCloser struct {
	items []ListenCloser
}

func (this CompositeListenCloser) Listen() {
	waiter := &sync.WaitGroup{}
	waiter.Add(len(this.items))

	for _, item := range this.items {
		this.beginListen(item, waiter)
	}

	waiter.Wait()
}
func (this CompositeListenCloser) beginListen(item Listener, waiter *sync.WaitGroup) {
	go func() {
		item.Listen()
		waiter.Done()
	}()
}

func (this CompositeListenCloser) Close() error {
	for _, item := range this.items {
		_ = item.Close()
	}

	return nil
}
