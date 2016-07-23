package goagrep

import (
	"fmt"
	"testing"
)

func BenchmarkGenerateDB(b *testing.B) {
	VERBOSE = false
	wordpath := "../example/testlist"
	path := "testlist.db"
	tupleLength := 4
	words, tuples, _, _ := scanWords(wordpath, tupleLength, false)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dumpToBoltDB(path, words, tuples, tupleLength)
	}
}

func BenchmarkGenerateDBInMemory(b *testing.B) {
	VERBOSE = false
	wordpath := "../example/testlist"
	tupleLength := 4
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GenerateDBInMemory(wordpath, tupleLength, VERBOSE)
	}
}

func BenchmarkScanWords(b *testing.B) {
	VERBOSE = false
	wordpath := "../example/testlist"
	tupleLength := 4
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		scanWords(wordpath, tupleLength, false)
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
