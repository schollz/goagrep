package goagrep

import (
	"fmt"
	"testing"
)

func init() {
	VERBOSE = false
}

func BenchmarkPartials(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getPartials("alligator", 3)
	}
}

func BenchmarkMatch(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GetMatch("heroint", "testlist3.db")
	}
}

func Example1() {
	wordpath := "../example/testlist"
	path := "testlist3.db"
	words, tuples := scanWords(wordpath, path, 3)
	dumpToBoltDB(path, words, tuples, 3)
	_, _, pairlist, _ := GetMatch("heroint", "testlist3.db")
	fmt.Println(pairlist[0])
	// Output: {heroine 1}
}

//
// func Example2() {
// 	_, _, pairlist, _ := GetMatch("zack's barn", "testlist3.db")
// 	fmt.Println(pairlist[0])
// 	// Output: {zack's barn 0}
// }
//
// func Example3() {
// 	GetMatch("zzzz zzzzzz", "testlist3.db")
// 	fmt.Println("got error")
// 	// Output: got error
// }
