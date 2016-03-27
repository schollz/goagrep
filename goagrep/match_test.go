package goagrep

import (
	"fmt"
	"testing"
)

func init() {
	VERBOSE = false
}

func Example1() {
	wordpath := "../example/testlist"
	path := "testlist3.db"
	words, tuples := scanWords(wordpath, path, 3)
	dumpToBoltDB(path, words, tuples, 3)
	match, score, err := GetMatch("heroint", "testlist3.db")
	fmt.Println(match, score, err)
	// Output: heroine 1 <nil>
}

func Example2() {
	wordpath := "../example/testlist"
	path := "testlist4.db"
	words, tuples := scanWords(wordpath, path, 4)
	dumpToBoltDB(path, words, tuples, 4)
	matches, scores, err := GetMatches("zack's barn", "testlist4.db")
	fmt.Println(matches[0:2], scores[0:2], err)
	// Output: [zack's barn zack's burn] [0 1] <nil>
}

func Example3() {
	wordpath := "../example/testlist"
	path := "testlist5.db"
	words, tuples := scanWords(wordpath, path, 5)
	dumpToBoltDB(path, words, tuples, 5)
	match, score, err := GetMatch("zzzzz zzz zzzz", "testlist5.db")
	fmt.Println(match, score, err)
	// Output: -1 No matches
}

func Example4() {
	matches, scores, err := GetMatches("zzzzz zzz zzzz", "testlist5.db")
	fmt.Println(matches, scores, err)
	// Output: [] [] No matches
}

func BenchmarkPartialsTuple3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getPartials("alligator", 3)
	}
}

func BenchmarkPartialsTuple4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getPartials("alligator", 4)
	}
}

func BenchmarkPartialsTuple5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getPartials("alligator", 5)
	}
}

func BenchmarkMatchTuple3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetMatch("heroint", "testlist3.db")
		GetMatch("myxovirus", "testlist3.db")
		GetMatch("pocket-handkerchief", "testlist3.db")
	}
}

func BenchmarkMatchTuple4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetMatch("heroint", "testlist4.db")
		GetMatch("myxovirus", "testlist4.db")
		GetMatch("pocket-handkerchief", "testlist4.db")
	}
}

func BenchmarkMatchTuple5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetMatch("heroint", "testlist5.db")
		GetMatch("pocket-handkerchief", "testlist5.db")
		GetMatch("myxovirus", "testlist5.db")
	}
}

func BenchmarkMatchesTuple5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetMatches("heroint", "testlist5.db")
		GetMatches("pocket-handkerchief", "testlist5.db")
		GetMatches("myxovirus", "testlist5.db")
	}
}
