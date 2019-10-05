package disruptor

import "sync"

type Disruptor struct {
	sequencer Sequencer
	workers   []ListenCloser
}

func (this Disruptor) Listen() {
	var waiter sync.WaitGroup
	waiter.Add(len(this.workers))
	this.listen(waiter)
	waiter.Wait()
}
func (this Disruptor) listen(waiter sync.WaitGroup) {
	for _, worker := range this.workers {
		go func(listener Listener) {
			listener.Listen()
			waiter.Done()
		}(worker)
	}
}

func (this Disruptor) Close() error {
	for _, reader := range this.workers {
		_ = reader.Close()
	}

	return nil
}

func (this Disruptor) Sequencer() Sequencer { return this.sequencer }
