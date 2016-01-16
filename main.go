package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "builddb") {
		tupleLength, _ := strconv.Atoi(os.Args[4])
		generateDB(strings.ToLower(os.Args[2]), strings.ToLower(os.Args[3]), tupleLength)
	} else if len(os.Args) > 1 {
		fmt.Println(getMatch(strings.ToLower(os.Args[2]), strings.ToLower(os.Args[1])))
	} else {
		fmt.Printf("Version 1.3")
		fmt.Println("Usage: ./main db 'word'")
	}
}
