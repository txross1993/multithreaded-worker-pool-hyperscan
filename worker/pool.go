package worker

import (
	"fmt"
	"sync"

	"github.com/flier/gohs/hyperscan"
)

type Pool struct {
	Tasks chan Task
	End   chan bool
	wg    *sync.WaitGroup
}

func StartPool(patterns []string, size int, handler hyperscan.MatchHandler) (*Pool, error) {

	bdb, err := newBlockDB(patterns)
	if err != nil {
		return nil, err
	}

	tasks := make(chan Task, size*10)
	end := make(chan bool)
	wg := &sync.WaitGroup{}
	pool := Pool{
		Tasks: tasks,
		End:   end,
		wg:    wg,
	}

	workers := make([]Worker, size)
	workerChan := make(chan chan Task, size)

	for i := 0; i < size; i++ {
		w, err := NewWorker(i, bdb, workerChan, handler)
		if err != nil {
			return nil, err
		}
		w.start()
		workers[i] = w
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-end:
				// drain final tasks; expects task chan to be closed and no more writes
				finalTasks := make([]Task, 0)
				for t := range tasks {
					finalTasks = append(finalTasks, t)
				}

				if len(finalTasks) > 0 {
					fmt.Printf("waiting for %d tasks to drain\n", len(finalTasks))
				}

				for _, task := range finalTasks {
					worker := <-workerChan
					worker <- task
				}

				for _, w := range workers {
					fmt.Println("stopping worker: ", w.id)
					w.stop()
				}
				return
			case task := <-tasks:
				worker := <-workerChan
				worker <- task
			}
		}
	}()

	return &pool, nil
}

func (p Pool) Stop() {
	p.End <- true
	p.wg.Wait()
}

func buildPatterns(patterns []string) []*hyperscan.Pattern {
	var compiled = make([]*hyperscan.Pattern, 0, len(patterns))

	for i, p := range patterns {
		word := fmt.Sprintf("\\b%s\\b", p)
		c := hyperscan.NewPattern(word, hyperscan.Caseless|hyperscan.DotAll|hyperscan.MultiLine|hyperscan.SingleMatch)
		c.Id = i
		compiled = append(compiled, c)
	}

	return compiled
}

func newBlockDB(patterns []string) (hyperscan.BlockDatabase, error) {
	hsPatterns := buildPatterns(patterns)
	return hyperscan.NewBlockDatabase(hsPatterns...)
}
