package workpool

import (
	"context"
	"sync"
)

type WorkPool struct {
	nWorkers int
}

func New(nWorkers int) WorkPool {
	return WorkPool{nWorkers: nWorkers}
}

func (wp *WorkPool) DoParallel(ctx context.Context, nJobs int, f func(workIndex int)) {
	queue := make(chan int)

	wg := sync.WaitGroup{}

	for i := 0; i < wp.nWorkers; i++ {
		go func() {
		loop:
			for {
				select {
				case wi, more := <-queue:
					if !more {
						break loop
					}

					defer wg.Done()
					f(wi)
				case <-ctx.Done():
					break loop
				}
			}

		}()
	}

	for i := 0; i < nJobs; i++ {
		wg.Add(1)
		select {
		case queue <- i:
		case <-ctx.Done():
			print("Canceled\n")
			return
		}
	}

	close(queue)

	wg.Wait()
}
