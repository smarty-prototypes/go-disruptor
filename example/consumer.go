package main

import "bitbucket.org/jonathanoliver/go-disruptor"

func consume(barrier disruptor.Barrier, source, sequence *disruptor.Sequence) {
	worker := disruptor.NewWorker(barrier, ConsumerHandler{}, source, sequence)

	for {
		worker.Process()
	}
}
