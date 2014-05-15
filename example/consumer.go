package main

import "bitbucket.org/jonathanoliver/go-disruptor"

func consume(barrier disruptor.Barrier, source, sequence *disruptor.Sequence) {
	worker := disruptor.NewWorker(barrier, TestHandler{}, source, sequence)

	for {
		worker.Process()
	}
}
