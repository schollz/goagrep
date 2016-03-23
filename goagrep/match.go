package goagrep

import (
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/arbovm/levenshtein"
	"github.com/boltdb/bolt"
)

var matches map[string]int

// GetMatch searches in the specified goagrep database.
// It returns the closest matched string and the Levenshtein distance.
//
// s is the string you want to search
//
// path is the filename of the database generated with GenerateDB()
func GetMatch(s string, path string) (string, int, PairList) {
	// normalize
	s = strings.ToLower(s)
	// Open a new bolt database
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tupleLength := 3
	wordBuckets := -1

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("vars"))
		v := b.Get([]byte("tupleLength"))
		tupleLength, _ = strconv.Atoi(string(v))
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("vars"))
		v := b.Get([]byte("wordBuckets"))
		wordBuckets, _ = strconv.Atoi(string(v))
		return nil
	})

	partials := getPartials(s, tupleLength)
	matches := make(map[string]int)
	// var wg sync.WaitGroup

	for _, partial := range partials {
		// wg.Add(1)
		func(partial string, path string) {
			// defer wg.Done()
			db.View(func(tx *bolt.Tx) error {
				var v []byte

				firstLetter := string(partial[0])
				secondLetter := string(partial[1])
				if strings.Contains(alphabet, firstLetter) && strings.Contains(alphabet, secondLetter) {
					b := tx.Bucket([]byte("tuples-" + firstLetter + secondLetter))
					v = b.Get([]byte(string(partial)))
				} else {
					b := tx.Bucket([]byte("tuples"))
					v = b.Get([]byte(string(partial)))
				}

				vals := string(v)
				// log.Println(partial)
				// log.Printf("The answer is: %v\n", vals)
				if len(v) > 0 {
					gotZero := false
					for _, k := range strings.Split(vals, " ") {
						db.View(func(tx *bolt.Tx) error {
							knum, _ := strconv.Atoi(k)
							b := tx.Bucket([]byte("words-" + strconv.Itoa(int(math.Mod(float64(knum), float64(wordBuckets))))))
							v := string(b.Get([]byte(k)))
							_, ok := matches[v]
							if ok != true {
								matches[v] = levenshtein.Distance(s, v)
								// fmt.Printf("Word match: %v\n", v)
								// fmt.Printf("Distance : %v\n", matches[v])
								if matches[v] == 0 {
									gotZero = true
								}
							}
							return nil
						})
						if gotZero {
							break
						}
					}
				}
				return nil
			})
		}(partial, path)
	}

	// wg.Wait()
	bestMatch := "none"
	bestVal := 100
	for k, v := range matches {
		if v < bestVal {
			bestMatch = k
			bestVal = v
		}

	}

	// Return at most 100 of the pairs
	pairlist := rankByWordCount(matches)
	if len(pairlist) > 100 {
		pairlist = pairlist[0:99]
	}
	return bestMatch, bestVal, pairlist
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
