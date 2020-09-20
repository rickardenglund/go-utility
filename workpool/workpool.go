package workpool

import "sync"

type WorkPool struct {
	nWorkers int
}

func New(nWorkers int) WorkPool {
	return WorkPool{nWorkers: nWorkers}
}

func (wp *WorkPool) DoParallel(nJobs int, f func(workIndex int)) {
	queue := make(chan int, wp.nWorkers)

	wg := sync.WaitGroup{}
	wg.Add(nJobs)

	for i := 0; i < wp.nWorkers; i++ {
		go func() {
			workIndex, more := <-queue
			if !more {
				return
			}

			defer wg.Done()
			f(workIndex)
		}()
	}

	for i := 0; i < nJobs; i++ {
		queue <- i
	}

	close(queue)

	wg.Wait()
}
