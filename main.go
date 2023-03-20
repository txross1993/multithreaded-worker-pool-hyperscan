package main

import (
	"fmt"
	"log"
	"runtime"
	"runtime/pprof"

	td "hyperscan-cgo/testdata"
	"hyperscan-cgo/worker"

	_ "go.uber.org/automaxprocs"
)

func main() {
	fmt.Println("num cpu: ", runtime.NumCPU())
	patterns := td.Patterns
	input := td.TestData
	expectedMatches := td.ExpectedMatches

	size := 8

	reader := worker.NewReader(size * 10)
	p, err := worker.StartPool(patterns, size, reader.Callback)
	if err != nil {
		log.Fatal(err)
	}

	for i, str := range input {
		p.Tasks <- worker.Task{
			ID:    i,
			Input: []byte(str),
		}
	}
	close(p.Tasks)
	p.Stop()
	reader.Stop()
	pprof.StopCPUProfile()
	got := reader.Matches()

	for i := 0; i < len(input); i++ {
		testCase := input[i]
		fmt.Println("running test case: ", testCase)

		actual, ok := got[i]
		if !ok && len(expectedMatches[i]) != 0 {
			fmt.Printf("\tDid not get expected match for test case %s\n", testCase)
			continue
		}
		expected := expectedMatches[i]

		if len(expected) != len(actual) {
			fmt.Printf("\tWANT %d matches GOT %d: %v\n", len(expected), len(actual), actual)
			continue
		}

		for j := 0; j < len(actual); j++ {
			if expected[j] != actual[j] {
				fmt.Printf("\tWANT matching pattern %d, GOT matching pattern %d\n", expected[j], actual[j])
			}
		}
	}

}
