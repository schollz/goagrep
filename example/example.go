package main

import (
	"fmt"
	"os"

	"github.com/schollz/goagrep/goagrep"
)

func init() {
	// Build database
	// only needs to be done once!
	databaseFile := "words.db"
	wordlist := "testlist"
	tupleLength := 3
	if _, err := os.Stat(databaseFile); os.IsNotExist(err) {
		goagrep.GenerateDB(wordlist, databaseFile, tupleLength)
	}

}
func main() {
	// Find word
	databaseFile := "words.db"
	searchWord := "heroint"
	word, score := goagrep.GetMatch(searchWord, databaseFile)
	fmt.Println(word, score)
}
