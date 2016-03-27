package goagrep

import "testing"

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
