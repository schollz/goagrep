# go-string-matching


Benchmarking (searching for 'madgascar' with a cache containing 3-letter tuples). Run 1000 times and taking average.

| Language | Runtime  |
|--------|--------|
| Python | 21 ms  |
| Go N=1 | 3-6 ms |
| Go N=8 | 2-4 ms |

# Use

```
go build match.go
```

To use, you first must build a database of words

```
./match build wordlist | sqlite3 words.db
```

Then to run simply use

```
./match "word or phrase"
```

# To do

- User friendliness?
- ~Command line help~
- ~Command line for generating cache~
- ~Convert to lowercase for converting~
