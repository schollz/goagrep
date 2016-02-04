package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/schollz/goagrep/goagrep"
)

var alphabet string

func main() {
	app := cli.NewApp()
	app.Name = "goagrep"
	app.Usage = "Fuzzy matching of big strings.\n   Before use, make sure to make a data file (go-agrep build help)."
	app.Version = "1.26"
	alphabet = "abcdefghijklmnopqrstuvwxyz-"
	var wordlist, subsetSize, outputFile, searchWord string

	app.Commands = []cli.Command{
		{
			Name:    "match",
			Aliases: []string{"m"},
			Usage:   "fuzzy match word",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "database, d",
					Usage:       "input database name (built using 'go-agrep build')",
					Destination: &wordlist,
				},
				cli.StringFlag{
					Name:        "word, w",
					Usage:       "word to use",
					Destination: &searchWord,
				},
			},
			Action: func(c *cli.Context) {
				if len(wordlist) == 0 || len(searchWord) == 0 {
					cli.ShowCommandHelp(c, "match")
				} else {
					word, score := goagrep.GetMatch(strings.ToLower(searchWord), wordlist)
					fmt.Printf("%v|||%v", word, score)
				}
			},
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "builds the database subsequent fuzzy matching",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "list, l",
					Usage:       "wordlist to use, seperated by newlines",
					Destination: &wordlist,
				},
				cli.StringFlag{
					Name:        "database, d",
					Usage:       "output database name (default: words.db)",
					Destination: &outputFile,
				},
				cli.StringFlag{
					Name:        "size, s",
					Usage:       "subset size (default: 3)",
					Destination: &subsetSize,
				},
			},
			Action: func(c *cli.Context) {
				if len(subsetSize) == 0 {
					subsetSize = "3"
				}
				if len(outputFile) == 0 {
					outputFile = "words.db"
				}
				if len(wordlist) == 0 {
					cli.ShowCommandHelp(c, "build")
				} else {
					fmt.Println("Generating '" + outputFile + "' from '" + wordlist + "' with subset size " + subsetSize)
					tupleLength, _ := strconv.Atoi(subsetSize)
					goagrep.GenerateDB(wordlist, outputFile, tupleLength)
					fmt.Println("Finished building db")
				}
			},
		},
	}

	app.Run(os.Args)
}
