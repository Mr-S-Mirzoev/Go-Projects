package main

import "sync"

func pipelineItem(process job, in, out chan interface{}, waiter *sync.WaitGroup) {
	defer close(out)
	defer waiter.Done()

	process(in, out)
}

func ExecutePipeline(fs ...job) {
	wg := new(sync.WaitGroup)
	wg.Add(len(fs))

	var previous chan interface{}

	for _, f := range fs {
		next := make(chan interface{})
		go pipelineItem(f, previous, next, wg)
		previous = next
	}

	wg.Wait()
}
