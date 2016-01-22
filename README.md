# Go-String-Matching version 1.1
I've written several apps that allow users to search a database for music artists, or book titles, or other lists of words. To make my life easy, I often use make the searched word a primary key in the database. However, this assumes (incorrectly) that a user spells their search word correctly and exactly how I have it in my database. To overcome this I wrote this [fuzzy string matching](https://en.wikipedia.org/wiki/Approximate_string_matching) program that simply takes any string, mispelled or not, and matches to one in my key list.

# Benchmark
Benchmarking using the 1000-word `testlist`, run with `go test -bench=.` using Intel(R) Core(TM) i7-3770 CPU @ 3.40GHz. The Python benchmark was run using the same words and the same subset length.

Version                                                                               | Runtime | Memory | Filesize
------------------------------------------------------------------------------------- | ------- | ------ | --------
[Python](https://github.com/schollz/string_matching)                                  | 104 ms  | ~30 MB | 140K
[Go Sqlite3](https://github.com/schollz/go-string-matching/tree/sqlite3)              | 6.2 ms  | ~20 MB | 124K
[Go BoltDB (this version)](https://github.com/schollz/go-string-matching/tree/master) | 2.8 ms  | ~14 MB | 512K

## How does it work

This program splits search-words into smaller subsets, and then finds the corresponding known words that contain each subset. It then runs Levenshtein's algorithm on the new list of known words to find the best match to the search-word. This *greatly* decreases the search space and thus increases the matching speed.

The subset length dictates how many pieces a word should be cut into, for purposes of finding partial matches for mispelled words. For instance example: a subset length of 3 for the word "olive" would index "oli", "liv", and "ive". This way, if one searched for "oliv" you could still return "olive" since subsets "oli" and "liv" can still grab the whole word and check its Levenshtein distance (which should be very close as its only missing the one letter). 

A smaller subset length will be more forgiving (it allows more mispellings), thus more accurate, but it would require more disk and more time to process since there are more words for each subset. A bigger subset length will help save hard drive space and decrease the runtime since there are fewer words that have the same, longer, subset. Here are some benchmarks of searching for various words of different lengths:

### Subset benchmarking

Tested using 69,905 words and the version with BoltDB (1.1) and  Intel(R) Core(TM) i7-3770 CPU @ 3.40GHz.

Subset length | Runtime   | Filesize
------ | --------- | --------
2      | 76 - 126 ms | 8MB
3      | 7 - 29 ms | 32MB
4      | 4 - 15 ms | 32MB
5      | 3 - 9 ms  | 32MB

These results show that you can get much faster speeds with shorter subset lengths, but keep in mind that this will not be able to match strings that have an error in the middle of the string and are have a length < 2*subset length - 1.

# Setup
Install dependencies

```bash
go get github.com/arbovm/levenshtein
go get github.com/boltdb/bolt
```

Build using

```bash
go build
```

# Run
To use, you first must build a database of words (here using a subset size of 3):

```
$ ./go-string-matching* testlist words.db 3
Finished building db
```

And then you can match any of the words using:

```
$ ./go-string-matching* words.db "heroes"
heroine|||3
```

which returns the best match and the levenshtein score.

You can test with a big list of words from Univ. Michigan:

```bash
wget http://www-personal.umich.edu/%7Ejlawler/wordlist
```

# To do
- Make commmand line stuff with github.com/codegangsta/cli
- ~Command line help~
- ~Command line for generating cache~
- ~Convert to lowercase for converting~
- Handle case that word definetly does not exist
