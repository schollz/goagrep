package goagrep

import (
	"fmt"
	"testing"
)

func BenchmarkScanWordsTuple4(b *testing.B) {
	VERBOSE = false
	wordpath := "../example/testlist"
	path := "testlist.db"
	tupleLength := 4
	words, tuples := scanWords(wordpath, path, tupleLength)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dumpToBoltDB(path, words, tuples, tupleLength)
	}
}

func BenchmarkSplitIntoPartialsTuple4(b *testing.B) {
	VERBOSE = false
	wordpath := "../example/testlist"
	path := "testlist.db"
	tupleLength := 4
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		scanWords(wordpath, path, tupleLength)
	}
}

func ExampleParse1() {
	fmt.Println(getPartials("hairbrush", 5))
	// Output: [hairb airbr irbru rbrus brush]
}

func ExampleParse2() {
	fmt.Println(getPartials("The Story and of Some/thing", 5))
	// Output: [thest hesto estor story torya oryan ryand yands andso ndsom dsome somet ometh methi ethin thing]
}

func ExampleParse3() {
	fmt.Println(getPartials("hi", 5))
	// Output: [zzz]
}
