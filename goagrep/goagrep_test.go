package goagrep

import "testing"
import "fmt"

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
		GetMatch("heroint", path)
	}
}

func Example1() {
	_, _, pairlist, _ := GetMatch("heroint", path)
	fmt.Println(pairlist[0])
	// Output: {heroine 1}
}

func Example2() {
	_, _, pairlist, _ := GetMatch("zack's barn", path)
	fmt.Println(pairlist[0])
	// Output: {zack's barn 0}
}

func Example3() {
	_, _, _, err := GetMatch("zzzz zzzzzz", path)
	if err != nil {
		fmt.Println("got error")
	}
	// Output: got error
}

// func BenchmarkDB(b *testing.B) {
// 	generateDB("testlist", "gotests.db", 3)
// }

func init() {
	VERBOSE = false
	wordpath = "../example/testlist"
	path = "testlist.db"
	tupleLength = 4
	words, tuples = scanWords(wordpath, path, tupleLength)
	dumpToBoltDB(path, words, tuples, tupleLength)
}
