package worker

import (
	"fmt"
	"runtime"

	"github.com/flier/gohs/hyperscan"
)

type Worker struct {
	id         int
	bdb        hyperscan.BlockDatabase
	scratch    *hyperscan.Scratch
	workerChan chan chan Task
	tasks      chan Task
	end        chan bool
	handler    hyperscan.MatchHandler
}

func NewWorker(id int, bdb hyperscan.BlockDatabase, workerChan chan chan Task, handler hyperscan.MatchHandler) (Worker, error) {
	scratch, err := hyperscan.NewScratch(bdb)
	if err != nil {
		return Worker{}, err
	}
	w := Worker{
		id:         id,
		bdb:        bdb,
		scratch:    scratch,
		tasks:      make(chan Task),
		end:        make(chan bool),
		handler:    handler,
		workerChan: workerChan,
	}
	return w, nil
}

func (w Worker) start() {
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		w.startNoLock()
		return
	}()
}

func (w Worker) startNoLock() {
	wc := w.workerChan
	for {
		select {
		case task := <-w.tasks:
			if err := w.bdb.Scan(task.Input, w.scratch, w.handler, task); err != nil {
				fmt.Println(fmt.Errorf("%w: %q, %v %v %v", err, string(task.Input), w.scratch, w.handler, task))
			}
		case wc <- w.tasks:
		case <-w.end:
			fmt.Printf("worker %d received signal to terminate\n", w.id)
			return
		}
	}
}

func (w Worker) stop() {
	w.end <- true
}
