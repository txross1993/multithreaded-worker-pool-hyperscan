package worker

import (
	td "hyperscan-cgo/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	numWorkers := 3
	reader := NewReader(len(td.TestData))
	p, err := StartPool(td.Patterns, numWorkers, reader.Callback)
	if err != nil {
		t.Fatal(err)
	}

	for _, task := range getTasks() {
		p.Tasks <- task
	}
	close(p.Tasks)
	p.Stop()
	reader.Stop()
	got := reader.Matches()
	verifyResultsT(t, got)

}

func verifyResultsT(t *testing.T, got map[int][]int) {
	for i := 0; i < len(td.TestData); i++ {
		testCase := td.TestData[i]

		actual, ok := got[i]
		if !ok && len(td.ExpectedMatches[i]) != 0 {
			t.Errorf("\tDid not get expected match for test case %s\n", testCase)
			continue
		}
		expected := td.ExpectedMatches[i]

		if len(expected) != len(actual) {
			t.Errorf("\tTestCase: %s\n\t WANT %d matches (%v) GOT %d (%v)\n", testCase, len(expected), expected, len(actual), actual)
			continue
		}

		for j := 0; j < len(actual); j++ {
			assert.Containsf(t, actual, expected[j], "\tWANT matching pattern %d, GOT matching pattern %d\n", expected[j], actual[j])
		}
	}
}
