# go-string-matching
Benchmarking (searching for 'madgascar' with a cache containing 3-letter tuples). Run 1000 times and taking average.

Language   | Runtime | Memory                                | Filesize
---------- | ------- | ------------------------------------- | --------
Python     | 21 ms   | Requires loading all data into memory | ?
Go Sqlite3 | 6.2 ms  | ~20 MB                                | 124K
Go BoltDB  | 2.8 ms  | ~14 MB                                | 512K

# Run
To use, you first must build a database of words:

```
wget http://www-personal.umich.edu/%7Ejlawler/wordlist
```

And then build the word list

```
./match build wordlist | sqlite3 words.db
```

Then to run simply use

```
./match "madgascar"
```

# To do
- Make commmand line stuff with github.com/codegangsta/cli
- ~Command line help~
- ~Command line for generating cache~
- ~Convert to lowercase for converting~

## Scratch
``` go build *.go rm test.db && time ./match build wordlist | sqlite3 test.db          1.734 time ./match 'test' test.db                             0.005

go build _.go && time ./main builddb wordlist my.db 3    3min19s go build _.go && time ./main my.db "test"         0.009

go test -bench=.

BenchmarkPartials-8         1000           1664637 ns/op BenchmarkBoltDump-8     2000000000               0.02 ns/op BenchmarkMatch-8             500           2924633 ns/op

``
