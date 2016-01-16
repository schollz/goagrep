# Go-String-Matching version 1.1

# Benchmark

Benchmarking using the 1000-word `testlist`, run with `go test -bench=.` using Intel(R) Core(TM) i7-3770 CPU @ 3.40GHz.

Version   | Runtime | Memory | Filesize
---------- | ------- | ------ | --------
[Python](https://github.com/schollz/string_matching)     | 104 ms  | ~30 MB | 140K
[Go Sqlite3](https://github.com/schollz/go-string-matching/tree/sqlite3) | 6.2 ms  | ~20 MB | 124K
[Go BoltDB (this version)](https://github.com/schollz/go-string-matching/tree/master)  | 2.8 ms  | ~14 MB | 512K

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
To use, you first must build a database of words (here using a tuple size of 3):

```
$ ./go-string-matching* testlist newdb 3
Finished building db
```

And then you can match any of the words using:

```
$ ./go-string-matching* newdb "heroes"
heroine|||3
```

which returns the best match and the levenshtein score.

# To do
- Make commmand line stuff with github.com/codegangsta/cli
- ~Command line help~
- ~Command line for generating cache~
- ~Convert to lowercase for converting~
