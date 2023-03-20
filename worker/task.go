package worker

type Task struct {
	ID    int
	Input []byte
}

type Match struct {
	ID      int
	Pattern int
}
