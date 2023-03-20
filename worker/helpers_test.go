package worker

import td "hyperscan-cgo/testdata"

func submitTasks(tasks chan Task) {
	for i := 0; i < len(td.TestData); i++ {
		data := []byte(td.TestData[i])
		tasks <- Task{
			ID:    i,
			Input: data,
		}
	}
}

func getTasks() []Task {
	var tasks = make([]Task, len(td.TestData))
	for i := 0; i < len(td.TestData); i++ {
		data := []byte(td.TestData[i])
		tasks[i] = Task{
			ID:    i,
			Input: data,
		}
	}
	return tasks
}
