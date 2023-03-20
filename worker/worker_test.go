package worker

import (
	"testing"

	td "hyperscan-cgo/testdata"
)

var result error

func BenchmarkWorkerLock(b *testing.B) {
	var r error

	for n := 0; n < b.N; n++ {
		r = benchmarkWorker(b)
	}

	result = r
}

func benchmarkWorker(b *testing.B) error {
	reader := NewReader(len(td.TestData))
	pool, err := StartPool(td.Patterns, 1, reader.Callback)
	if err != nil {
		b.Fatal(err)
	}

	submitTasks(pool.Tasks)
	close(pool.Tasks)
	pool.Stop()
	reader.Stop()
	verifyResults(b, reader.Matches())

	return nil
}

func verifyResults(b *testing.B, got map[int][]int) {
	for i := 0; i < len(td.TestData); i++ {
		testCase := td.TestData[i]

		actual, ok := got[i]
		if !ok && len(td.ExpectedMatches[i]) != 0 {
			b.Errorf("\tDid not get expected match for test case %s\n", testCase)
			continue
		}
		expected := td.ExpectedMatches[i]

		if len(expected) != len(actual) {
			b.Errorf("\tTestCase: %s\n\t WANT %d matches (%v) GOT %d (%v)\n", testCase, len(expected), expected, len(actual), actual)
			continue
		}

		for j := 0; j < len(actual); j++ {
			if expected[j] != actual[j] {
				b.Errorf("\tWANT matching pattern %d, GOT matching pattern %d\n", expected[j], actual[j])
			}
		}
	}
}
