![Version 1.2](https://img.shields.io/badge/version-1.2-brightgreen.svg?version=flat-square)

# fmbs - Fuzzy matching of big strings
_A simple program to do fuzzy matching for strings of any length._

![Big Fuzz Mascot](http://ecx.images-amazon.com/images/I/417W-2NwzpL._SX355_.jpg)

There are situations where you want to take the user's input and match a primary key in a database (for example, I ran into this in my web apps for finding [music artists](http://www.musicsuggestions.ninja/), and [book titles](http://booksuggestions.ninja/)). But, immediately a problem is introduced: _what happens if the user spells the primary key incorrectly?_ This [fuzzy string matching](https://en.wikipedia.org/wiki/Approximate_string_matching) program solves this problem - it takes any string, misspelled or not, and matches to one a specified key list.

# Benchmark
Benchmarking using the 1000-word `testlist`, run with `go test -bench=.` using Intel(R) Core(TM) i7-3770 CPU @ 3.40GHz. The Python benchmark was run using the same words and the same subset length. The current GNU standard, `agrep`, was benchmarked using [perf](http://askubuntu.com/questions/50145/how-to-install-perf-monitoring-tool/306683): `perf stat -r 500 -d agrep -By "heroint" testlist`.

Version                                                                 | Runtime | Memory | Database size
----------------------------------------------------------------------- | ------- | ------ | -------------------------------------
[Python](https://github.com/schollz/string_matching)                    | 104 ms  | ~30 MB | 140K
[Go Sqlite3](https://github.com/schollz/fmbs/tree/sqlite3)              | 6 ms    | ~20 MB | 124K
[Go BoltDB (this version)](https://github.com/schollz/fmbs/tree/master) | 2 ms    | <1 MB  | 256K
[agrep](https://en.wikipedia.org/wiki/Agrep)                            | 2 ms    | <1 MB  | 0 (no precomputed database nessecary)

## Why use `fmbs` over `agrep`?
It seems that `agrep` really a comparable choice for most applications. It does not require any database and its comparable speed to `fmbs`. However, there are situations where `fmbs` is more useful:
1. `fmbs` can search much longer strings: `agrep` is limited to 32 characters while `fmbs` is only limited to 500.
2. `fmbs` can handle more mistakes in a string: `agrep` is limited to edit distances of 8, while `fmbs` has no limit.
3. `fmbs` can be 3-10x faster: You can set higher subset lengths to get faster speeds than `agrep`. For example, in a list of 255,615 book names + authors, `agrep` took ~150 ms while `fmbs` took 8 - 40 ms (using a subset length of 5).

# How does it work
This program splits search-words into smaller subsets, and then finds the corresponding known words that contain each subset. It then runs Levenshtein's algorithm on the new list of known words to find the best match to the search-word. This _greatly_ decreases the search space and thus increases the matching speed.

The subset length dictates how many pieces a word should be cut into, for purposes of finding partial matches for mispelled words. For instance example: a subset length of 3 for the word "olive" would index "oli", "liv", and "ive". This way, if one searched for "oliv" you could still return "olive" since subsets "oli" and "liv" can still grab the whole word and check its Levenshtein distance (which should be very close as its only missing the one letter).

A smaller subset length will be more forgiving (it allows more mispellings), thus more accurate, but it would require more disk and more time to process since there are more words for each subset. A bigger subset length will help save hard drive space and decrease the runtime since there are fewer words that have the same, longer, subset. Here are some benchmarks of searching for various words of different lengths:

### Subset benchmarking
Tested using 69,905 words and the version with BoltDB (1.1) and  Intel(R) Core(TM) i7-3770 CPU @ 3.40GHz.

Subset length | Runtime     | Filesize
------------- | ----------- | --------
2             | 76 - 126 ms | 8MB
3             | 7 - 29 ms   | 32MB
4             | 4 - 15 ms   | 32MB
5             | 3 - 9 ms    | 32MB

These results show that you can get much faster speeds with shorter subset lengths, but keep in mind that this will not be able to match strings that have an error in the middle of the string and are have a length < 2*subset length - 1.

# Setup
## Build ...
Install dependencies

```bash
go get github.com/arbovm/levenshtein
go get github.com/boltdb/bolt
go get github.com/codegangsta/cli
```

Build using

```bash
go build
```

## ... or Install
Install using

```bash
go get -u github.com/schollz/fmbs
```

# Run
## Basic Usage
Building DB:

```
USAGE:
   fmbs build [command options] [arguments...]

OPTIONS:
   --list, -l           wordlist to use, seperated by newlines
   --database, -d       output database name (default: words.db)
   --size, -s           subset size (default: 3)
```

Matching:

```
USAGE:
   fmbs match [command options] [arguments...]

OPTIONS:
   --database, -d       input database name (built using 'fmbs build')
   --word, -w           word to use
```

## Example
First compile a list of your phrases or words that you want to match (see `testlist`). Then you can build a `fmbs` database using:

```
$ fmbs build -l testlist -d words.db
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
$ fmbs match -w heroint -d words.db
heroine|||1
```

which returns the best match and the levenshtein score.

You can test with a big list of words from Univ. Michigan:

```bash
wget http://www-personal.umich.edu/%7Ejlawler/wordlist
```

# To do
- ~~Make commmand line stuff with github.com/codegangsta/cli~~
- ~~Command line help~~
- ~~Command line for generating cache~~
- ~~Convert to lowercase for converting~~
- ~~Vastly increased DB building by decreasing size of partials (`make([]string, 500)`) and making extra buckets~~
- Handle case that word definetly does not exist
- Save searches, so caching can be used to find common searches easily
- Use channels for faster searching?
