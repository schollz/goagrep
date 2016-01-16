package main

import "testing"

var words map[string]int
var tuples map[string]string
var wordpath string
var path string
var tupleLength int

func BenchmarkPartials(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getPartials("alligator", 3)
	}
}

func BenchmarkBoltDump(b *testing.B) {
	dumpToBoltDB(path, words, tuples, tupleLength)
}

func BenchmarkMatch(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getMatch("heroint", path)
	}
}

// func BenchmarkDB(b *testing.B) {
// 	generateDB("testlist", "gotests.db", 3)
// }

func init() {
	wordpath = "testlist"
	path = "testlist.db"
	tupleLength = 3
	words, tuples = scanWords(wordpath, path, tupleLength)
	dumpToBoltDB(path, words, tuples, tupleLength)
}
