# go-string-matching

Benchmarking (searching for 'madgascar' with a cache containing 3-letter tuples). Run 1000 times and taking average.

| Language | Runtime  | Memory |
|--------|--------|--------|
| Python | 21 ms  | Requires loading all data into memory |
| Go | 15-20 ms | Requires loading no data into memory! |

# Install

First make sure you have Go installed then use:

```
go get github.com/arbovm/levenshtein
go get github.com/cheggaaa/pb
go get github.com/mattn/go-sqlite3
go get github.com/codegangsta/cli
go build match.go
```


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
