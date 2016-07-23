![Version 2.0beta](https://img.shields.io/badge/version-2.0beta-brightgreen.svg?version=flat-square) ![Coverage](https://img.shields.io/badge/coverage-78%25-orange.svg) [![GoDoc](https://godoc.org/github.com/schollz/goagrep/goagrep?status.svg)](https://godoc.org/github.com/schollz/goagrep/goagrep)

# goagrep

<!-- ![Big Fuzz Mascot](http://ecx.images-amazon.com/images/I/417W-2NwzpL._SX355_.jpg) -->

 _There are situations where you want to take the user's input and match a primary key in a database. But, immediately a problem is introduced: what happens if the user spells the primary key incorrectly? This fuzzy string matching program solves this problem - it takes any string, misspelled or not, and matches to one a specified key list._

# About

`goagrep` requires building a precomputed database from the file that has the target strings. Then, when querying, `goagrep` splits the search string into smaller subsets, and then finds the corresponding known target strings that contain each subset. It then runs Levenshtein's algorithm on the new list of target strings to find the best match to the search string. This _greatly_ decreases the search space and thus increases the matching speed.

The subset length dictates how many pieces a word should be cut into, for purposes of finding partial matches for mispelled words. For instance example: a subset length of 3 for the word "olive" would index "oli", "liv", and "ive". This way, if one searched for "oliv" you could still return "olive" since subsets "oli" and "liv" can still grab the whole word and check its Levenshtein distance (which should be very close as its only missing the one letter).

A smaller subset length will be more forgiving (it allows more mispellings), thus more accurate, but it would require more disk and more time to process since there are more words for each subset. A bigger subset length will help save hard drive space and decrease the runtime since there are fewer words that have the same, longer, subset. You can get much faster speeds with longer subset lengths, but keep in mind that this will not be able to match strings that have an error in the middle of the string and are have a length < 2*subset length - 1.

## Why use `goagrep`?

It seems that [`agrep`](https://github.com/Wikinaut/agrep) really a comparable choice for most applications. It does not require any database and its comparable speed to `goagrep`. However, there are situations where `goagrep` is more useful:

1. `goagrep` can search much longer strings: [`agrep`](https://github.com/Wikinaut/agrep) is limited to 32 characters while `goagrep` is only limited to 500\. [`tre-agrep`](http://laurikari.net/tre/download/) is not limited, but is much slower.
2. `goagrep` can handle many mistakes in a string: [`agrep`](https://github.com/Wikinaut/agrep) is limited to edit distances of 8, while `goagrep` has no limit.
3. `goagrep` is _fast_ (see benchmarking below), and the speed can be tuned: You can set higher subset lengths to get faster speeds and less accuracy - leaving the tradeoff up to you

## Benchmarking

Benchmarking using the [319,378 word dictionary](http://www.md5this.com/tools/wordlists.html) (3.5 MB), run with `perf stat -r 50 -d <CMD>` or using `go test -bench=Match` using AMD FX(tm)-8350.

Program                                         | Runtime | Memory usage
----------------------------------------------- | ------- | ------------
`goagrep` `in memory`, subset size = 5     | 0.2 ms  | 90 MB ram
`goagrep` `DB`, subset size = 5            | 0.9 ms    | 64 MB disk
`goagrep` `in memory`, subset size = 3     | 18 ms   | 90 MB ram
`goagrep` `DB`, subset size = 3            | 71 ms   | 64 MB disk
[agrep](https://github.com/Wikinaut/agrep)      | 7 ms    | 3.5 MB disk
[tre-agrep](http://laurikari.net/tre/download/) | 613 ms  | 3.5 MB disk

# Installation

```bash
go get github.com/schollz/goagrep
```

# Usage

You can either build a hard-disk database or use it in memory (i.e. in a program or as a TCP client). See `main.go` or tests for examples.

## TCP server example

Start a server with a list:

```
> $GOPATH/bin/goagrep serve -l ../example/testlist
2016/07/23 07:05:35 Creating server with address localhost:9992
```

And then you can match words using netcat or similar:

```
> echo heroint | nc localhost 9992
heroine
```


## Standalone cli utility

Build a database with a list:

```
> $GOPATH/bin/goagrep build -l ../example/testlist
Generating 'words.db' from '../example/testlist' with subset size 3
Finished building db
```

And then you can match by telling the program where the database is:

```
> $GOPATH/bin/goagrep match -d words.db -w heroint
heroine|||1
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
