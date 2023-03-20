package worker

import (
	"runtime"
)

type Reader struct {
	end     chan bool
	results chan Match
	matches map[int][]int
}

func NewReader(bufSize int) *Reader {
	r := &Reader{
		end:     make(chan bool),
		results: make(chan Match, bufSize),
		matches: make(map[int][]int),
	}

	go r.start()

	return r
}

func (r *Reader) start() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	for {
		select {
		case <-r.end:
			return
		case match := <-r.results:
			r.matches[match.ID] = append(r.matches[match.ID], match.Pattern)
		}
	}
}

func (r *Reader) Stop() {
	r.end <- true
	close(r.results)
	for match := range r.results {
		r.matches[match.ID] = append(r.matches[match.ID], match.Pattern)
	}
}

func (r *Reader) Matches() map[int][]int {
	return r.matches
}

func (r *Reader) Clear() {
	r.matches = make(map[int][]int)
}

func (r *Reader) Callback(id uint, _, _ uint64, _ uint, callbackCtx interface{}) error {
	task, ok := callbackCtx.(Task)
	if ok {
		r.results <- Match{
			ID:      task.ID,
			Pattern: int(id),
		}
	}
	// fmt.Println("no comprendo")
	// return errors.New("failed to understand context!")
	return nil
}
