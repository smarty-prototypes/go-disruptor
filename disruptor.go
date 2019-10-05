package disruptor

import "sync"

type Disruptor struct {
	sequencer Sequencer
	listeners []ListenCloser
}

func (this Disruptor) Listen() {
	waiter := &sync.WaitGroup{}
	waiter.Add(len(this.listeners))
	this.listen(waiter)
	waiter.Wait()
}
func (this Disruptor) listen(waiter *sync.WaitGroup) {
	for _, worker := range this.listeners {
		go func(listener Listener) {
			listener.Listen()
			waiter.Done()
		}(worker)
	}
}

func (this Disruptor) Close() error {
	for _, listener := range this.listeners {
		_ = listener.Close()
	}

	return nil
}

func (this Disruptor) Sequencer() Sequencer { return this.sequencer }
