package main

import (
	"fmt"
	"os"

	"github.com/schollz/goagrep/goagrep"
)

var databaseFile string
var wordlist string
var tupleLength int

func init() {
	databaseFile = "words.db"
	wordlist = "testlist"
	tupleLength = 5

	// Build database
	if _, err := os.Stat(databaseFile); os.IsNotExist(err) {
		goagrep.GenerateDB(wordlist, databaseFile, tupleLength, true)
	}
}

func main() {
	// Find word
	searchWord := "heroint"
	word, score, err := goagrep.GetMatch(searchWord, databaseFile)
	fmt.Println(word, score, err)
}
