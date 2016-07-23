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
	words, tuples, _, _ := scanWords(wordpath, 3, false)
	dumpToBoltDB(path, words, tuples, 3)
	match, score, err := GetMatch("heroint", "testlist3.db")
	fmt.Println(match, score, err)
	// Output: heroine 1 <nil>
}

func Example2() {
	wordpath := "../example/testlist"
	path := "testlist4.db"
	words, tuples, _, _ := scanWords(wordpath, 4, false)
	dumpToBoltDB(path, words, tuples, 4)
	matches, scores, err := GetMatches("zack's barn", "testlist4.db")
	fmt.Println(matches[0:2], scores[0:2], err)
	// Output: [zack's barn zack's burn] [0 1] <nil>
}

func Example3() {
	wordpath := "../example/testlist"
	path := "testlist5.db"
	words, tuples, _, _ := scanWords(wordpath, 5, false)
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

func Example5() {
	matches, err := findMatches("cambium", "testlist4.db")
	fmt.Println(len(matches), err)
	// Output: 3 <nil>
}

func Example6() {
	stringListPath := "../example/testlist"
	tupleLength := 3
	_, _, words, tuples := scanWords(stringListPath, tupleLength, true)
	matches, _, _ := GetMatchesInMemory("mykovirus", words, tuples, tupleLength, true)
	fmt.Println(matches[0])
	// Output: myxovirus
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

func BenchmarkMatchTuple3InMemory(b *testing.B) {
	stringListPath := "testlist"
	tupleLength := 3
	_, _, words, tuples := scanWords(stringListPath, tupleLength, true)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GetMatchesInMemory("heroint", words, tuples, tupleLength, true)
		GetMatchesInMemory("myxovirus", words, tuples, tupleLength, true)
		GetMatchesInMemory("pocket-handkerchief", words, tuples, tupleLength, true)
	}
}

func BenchmarkMatchTuple4InMemory(b *testing.B) {
	stringListPath := "../example/testlist"
	tupleLength := 4
	_, _, words, tuples := scanWords(stringListPath, tupleLength, true)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GetMatchesInMemory("heroint", words, tuples, tupleLength, true)
		GetMatchesInMemory("myxovirus", words, tuples, tupleLength, true)
		GetMatchesInMemory("pocket-handkerchief", words, tuples, tupleLength, true)
	}
}

func BenchmarkMatchTuple5InMemory(b *testing.B) {
	stringListPath := "testlist"
	tupleLength := 5
	_, _, words, tuples := scanWords(stringListPath, tupleLength, true)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GetMatchesInMemory("heroint", words, tuples, tupleLength, true)
		GetMatchesInMemory("myxovirus", words, tuples, tupleLength, true)
		GetMatchesInMemory("pocket-handkerchief", words, tuples, tupleLength, true)
	}
}
func BenchmarkMatchTuple3(b *testing.B) {
	wordpath := "testlist"
	path := "testlist3.db"
	words, tuples, _, _ := scanWords(wordpath, 3, false)
	dumpToBoltDB(path, words, tuples, 3)
	b.ResetTimer()
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
	wordpath := "testlist"
	path := "testlist5.db"
	words, tuples, _, _ := scanWords(wordpath, 5, false)
	dumpToBoltDB(path, words, tuples, 5)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GetMatch("heroint", "testlist5.db")
		GetMatch("pocket-handkerchief", "testlist5.db")
		GetMatch("myxovirus", "testlist5.db")
	}
}
