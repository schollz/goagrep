package goagrep

import (
	"errors"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/arbovm/levenshtein"
	"github.com/boltdb/bolt"
	"github.com/kmulvey/gohamming/hamming"
)

func GetMatchesInMemory(s string, wordsLookup map[int]string, tuplesLookup map[string][]int, tupleLength int, findBestMatch bool) ([]string, []int, error) {
	bestMatch := "ajcoewiclaksmecoiawemcolwqiemjclaseflkajsfklj"
	bestVal := 1000
	var returnError error
	returnError = nil
	s = strings.ToLower(s)
	partials := getPartials(s, tupleLength)
	matches := make(map[string]int)
	for _, partial := range partials {
		for _, val := range tuplesLookup[partial] {
			possibleWord := wordsLookup[val]
			if _, ok := matches[possibleWord]; !ok {
				matches[possibleWord] = hamming.Calc(s, strings.ToLower(possibleWord))
				if matches[possibleWord] < 0 {
					matches[possibleWord] = levenshtein.Distance(s, strings.ToLower(possibleWord))
				}
				if matches[possibleWord] < bestVal {
					bestMatch = possibleWord
					bestVal = matches[possibleWord]
				}
			}
		}
	}

	if findBestMatch {
		if bestMatch == "ajcoewiclaksmecoiawemcolwqiemjclaseflkajsfklj" {
			returnError = errors.New("No matches")
			bestMatch = ""
			bestVal = -1
		}
		return append([]string{}, bestMatch), append([]int{}, bestVal), nil
	}

	matchWords := []string{}
	matchScores := []int{}
	var pairlist PairList
	if len(matches) > 1 {
		pairlist = rankByWordCount(matches)
		if len(pairlist) > 100 {
			pairlist = pairlist[0:99]
		}
		for i := range pairlist {
			matchWords = append(matchWords, pairlist[i].Key)
			matchScores = append(matchScores, pairlist[i].Value)
		}
	} else {
		returnError = errors.New("No matches")
	}
	return matchWords, matchScores, returnError
}

func findMatches(s string, path string, bestMatchOnly bool) (matches map[string]int, bestMatch string, bestVal int, returnError error) {
	bestMatch = "ajcoewiclaksmecoiawemcolwqiemjclaseflkajsfklj"
	bestVal = 1000
	returnError = nil
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
	matches = make(map[string]int)
	// var wg sync.WaitGroup
	foundBestMatch := false
	for _, partial := range partials {
		if foundBestMatch && bestMatchOnly {
			break
		}
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
					for _, k := range strings.Split(vals, " ") {
						db.View(func(tx *bolt.Tx) error {
							knum, _ := strconv.Atoi(k)
							b := tx.Bucket([]byte("words-" + strconv.Itoa(int(math.Mod(float64(knum), float64(wordBuckets))))))
							v := string(b.Get([]byte(k)))
							_, ok := matches[v]
							if ok != true {
								matches[v] = hamming.Calc(s, strings.ToLower(v))
								if matches[v] < 0 {
									matches[v] = levenshtein.Distance(s, strings.ToLower(v))
								}
								if matches[v] < bestVal {
									bestMatch = v
									bestVal = matches[v]
									if bestVal == 0 {
										foundBestMatch = true
									}
								}
							}
							return nil
						})
						if foundBestMatch && bestMatchOnly {
							break
						}
					}
				}
				return nil
			})
		}(partial, path)
	}
	return
}

// GetMatch searches in the specified goagrep database.
// It returns the closest matched string and the Levenshtein distance.
//
// s is the string you want to search
//
// path is the filename of the database generated with GenerateDB()
func GetMatch(s string, path string) (string, int, error) {
	_, bestMatch, bestVal, returnError := findMatches(s, path, true)
	if bestMatch == "ajcoewiclaksmecoiawemcolwqiemjclaseflkajsfklj" {
		returnError = errors.New("No matches")
		bestMatch = ""
		bestVal = -1
	}
	return bestMatch, bestVal, returnError
}

// GetMatches searches in the specified goagrep database.
// Returns the a list of the at most 100 words and scores in order.
//
// s is the string you want to search
//
// path is the filename of the database generated with GenerateDB()
func GetMatches(s string, path string) ([]string, []int, error) {
	matches, _, _, returnError := findMatches(s, path, false)
	matchWords := []string{}
	matchScores := []int{}
	var pairlist PairList
	if len(matches) > 1 {
		pairlist = rankByWordCount(matches)
		if len(pairlist) > 100 {
			pairlist = pairlist[0:99]
		}
		for i := range pairlist {
			matchWords = append(matchWords, pairlist[i].Key)
			matchScores = append(matchScores, pairlist[i].Value)
		}
	} else {
		returnError = errors.New("No matches")
	}
	return matchWords, matchScores, returnError
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

// Pair structure for sorting
type Pair struct {
	Key   string
	Value int
}

// PairList array structure for sorting
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
