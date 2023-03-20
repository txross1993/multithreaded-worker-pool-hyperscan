package testdata

var Patterns = []string{
	"foo",
	"bar",
	"baz",
	"magic",
	"lemon",
}

var TestData = []string{
	"my magic lemon",
	"is a demon",
	"surprise this string does not match anything",
	"barney baz",
	"foo foo foo",
}

var ExpectedMatches = map[int][]int{
	0: {3, 4},
	3: {2},
	4: {0},
}
