package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) == 4 {
		tupleLength, _ := strconv.Atoi(os.Args[3])
		generateDB(strings.ToLower(os.Args[1]), strings.ToLower(os.Args[2]), tupleLength)
		fmt.Println("Finished building db")
	} else if len(os.Args) == 3 {
		word, score := getMatch(strings.ToLower(os.Args[2]), strings.ToLower(os.Args[1]))
		fmt.Printf("%v|||%v", word, score)
	} else {
		fmt.Printf("Version 1.1\n\n")
		fmt.Println(`Usage:

BUILDING THE DB:

./go-string-matching wordlist newdb numTuples

wordlist = a file with a list of words/phrases
newdb = the output processed boltdb
numTuples = number of tuples you want processed


MATCHING A STRING:

./go-string-matching newdb "string"

newdb = the output processed boltdb
"string" = what you want searched

`)
	}
}
