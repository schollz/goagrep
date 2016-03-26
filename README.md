![Version 1.2](https://img.shields.io/badge/version-1.2-brightgreen.svg?version=flat-square) [![GoDoc](https://godoc.org/github.com/schollz/goagrep/goagrep?status.svg)](https://godoc.org/github.com/schollz/goagrep/goagrep)

# goagrep

<!-- ![Big Fuzz Mascot](http://ecx.images-amazon.com/images/I/417W-2NwzpL._SX355_.jpg)
 -->

_There are situations where you want to take the user's input and match a primary key in a database. But, immediately a problem is introduced: what happens if the user spells the primary key incorrectly? This fuzzy string matching program solves this problem - it takes any string, misspelled or not, and matches to one a specified key list._

# About

`goagrep` requires building a precomputed database from the file that has the target strings. Then, when querying, `goagrep` splits the search string into smaller subsets, and then finds the corresponding known target strings that contain each subset. It then runs Levenshtein's algorithm on the new list of target strings to find the best match to the search string. This _greatly_ decreases the search space and thus increases the matching speed.

The subset length dictates how many pieces a word should be cut into, for purposes of finding partial matches for mispelled words. For instance example: a subset length of 3 for the word "olive" would index "oli", "liv", and "ive". This way, if one searched for "oliv" you could still return "olive" since subsets "oli" and "liv" can still grab the whole word and check its Levenshtein distance (which should be very close as its only missing the one letter).

A smaller subset length will be more forgiving (it allows more mispellings), thus more accurate, but it would require more disk and more time to process since there are more words for each subset. A bigger subset length will help save hard drive space and decrease the runtime since there are fewer words that have the same, longer, subset. You can get much faster speeds with longer subset lengths, but keep in mind that this will not be able to match strings that have an error in the middle of the string and are have a length < 2*subset length - 1.

## Why use `goagrep`?
It seems that [`agrep`](https://github.com/Wikinaut/agrep)  really a comparable choice for most applications. It does not require any database and its comparable speed to `goagrep`. However, there are situations where `goagrep` is more useful:

1. `goagrep` can search much longer strings: [`agrep`](https://github.com/Wikinaut/agrep)  is limited to 32 characters while `goagrep` is only limited to 500. [`tre-agrep`](http://laurikari.net/tre/download/)  is not limited, but is much slower.
2. `goagrep` can handle many mistakes in a string: [`agrep`](https://github.com/Wikinaut/agrep)  is limited to edit distances of 8, while `goagrep` has no limit.
3. `goagrep` is *fast* (see benchmarking below), and the speed can be tuned: You can set higher subset lengths to get faster speeds and less accuracy - leaving the tradeoff up to you

## Benchmarking
Benchmarking using the [319,378 word dictionary](http://www.md5this.com/tools/wordlists.html) (3.5 MB), run with `perf stat -r 50 -d <CMD>` using Intel(R) Core(TM) i5-4310U CPU @ 2.00GHz. These benchmarks were run with a single word, and can flucuate ~50% depending on the word.

Program                                             | Runtime  | Database size
--------------------------------------------------- | -------- | -----------------------
[goagrep](https://github.com/schollz/goagrep/tree/master) | **3 ms**     | 69 MB (subset size = 5)
[goagrep](https://github.com/schollz/goagrep/tree/master) | 7 ms     | 66 MB (subset size = 4)
[goagrep](https://github.com/schollz/goagrep/tree/master) | 84 ms    | 58 MB (subset size = 3)
[agrep](https://github.com/Wikinaut/agrep)          | 53 ms    | 3.5 MB (original file)
[tre-agrep](http://laurikari.net/tre/download/)     | 1,178 ms | 3.5 MB (original file)



# Installation

```bash
go get github.com/schollz/goagrep
```

# Usage (program)

## Building DB

```
USAGE:
   goagrep build [command options] [arguments...]

OPTIONS:
   --list, -l           wordlist to use, seperated by newlines
   --database, -d       output database name (default: words.db)
   --size, -s           subset size (default: 3)
```

## Matching

```
USAGE:
   goagrep match [command options] [arguments...]

OPTIONS:
   --database, -d       input database name (built using 'goagrep build')
   --word, -w           word to use
```

## Example
First compile a list of your phrases or words that you want to match (see `testlist`). Then you can build a `goagrep` database using:

```
$ goagrep build -l testlist -d words.db
Generating 'words.db' from 'testlist' with subset size 3
Parsing subsets...
1000 / 1000  100.00 % 0
Finished parsing subsets
Loading words into db...
1000 / 1000  100.00 % 0
Words took 13.0398ms
Loading subsets into db...
2281 / 2281  100.00 % 0
Subsets took 19.0267ms
Finished building db
```

And then you can match any of the words using:

```
$ goagrep match -w heroint -d words.db
heroine|||1
```

which returns the best match and the levenshtein score.

You can test with a big list of words from Univ. Michigan:

```bash
wget http://www-personal.umich.edu/%7Ejlawler/wordlist
```

# Usage (library)

You can also use as a library. Here's an example program (see in `example/`)

```golang
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
	word, score, _, _ := goagrep.GetMatch(searchWord, databaseFile)
	fmt.Println(word, score)
}
```

# History
- ~~Make commmand line stuff with github.com/codegangsta/cli~~
- ~~Command line help~~
- ~~Command line for generating cache~~
- ~~Convert to lowercase for converting~~
- ~~Vastly increased DB building by decreasing size of partials (`make([]string, 500)`) and making extra buckets~~
- Handle case that word definetly does not exist
- Save searches, so caching can be used to find common searches easily
- Use channels for faster searching?
- Add in `agrep` options
