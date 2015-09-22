# go-string-matching

Benchmarking (searching for 'madgascar' with a cache containing 3-letter tuples). Run 1000 times and taking average.

| Language | Runtime  |
|--------|--------|
| Python | 21 ms  |
| Go N=1 | 3-6 ms |
| Go N=8 | 2-4 ms |

# Install

First make sure you have Go installed then use:

```
go get github.com/arbovm/levenshtein
go get github.com/cheggaaa/pb
go get github.com/mattn/go-sqlite3
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

- User friendliness?
- ~Command line help~
- ~Command line for generating cache~
- ~Convert to lowercase for converting~
