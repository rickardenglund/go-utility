package workpool

import (
	"sync"
)

type WorkPool struct {
	nWorkers int
}

func New(nWorkers int) WorkPool {
	return WorkPool{nWorkers: nWorkers}
}

func (wp *WorkPool) DoParallel(nJobs int, f func(workIndex int)) {
	queue := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(nJobs)

	for i := 0; i < wp.nWorkers; i++ {
		go func() {
			for {
				workIndex, more := <-queue
				if !more {
					break
				}

				defer wg.Done()
				f(workIndex)
			}
		}()
	}

	for i := 0; i < nJobs; i++ {
		queue <- i
	}

	close(queue)

	wg.Wait()
}
