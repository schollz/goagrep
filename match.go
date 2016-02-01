package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/arbovm/levenshtein"
	"github.com/boltdb/bolt"
)

var matches map[string]int

func getMatch(s string, path string) (string, int) {
	// normalize
	s = strings.ToLower(s)
	// Open a new bolt database
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tupleLength := 3

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("vars"))
		v := b.Get([]byte("tupleLength"))
		tupleLength, _ = strconv.Atoi(string(v))
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
				if string(partial[0]) <= "c" { // DIVIDED 6x: 32MB 84 ms...UNDIVIDED: 188ms
					b1 := tx.Bucket([]byte("tuples-1"))
					v = b1.Get([]byte(string(partial)))
				} else if string(partial[0]) <= "f" {
					b2 := tx.Bucket([]byte("tuples-2"))
					v = b2.Get([]byte(string(partial)))
				} else if string(partial[0]) <= "i" {
					b3 := tx.Bucket([]byte("tuples-3"))
					v = b3.Get([]byte(string(partial)))
				} else if string(partial[0]) <= "l" {
					b4 := tx.Bucket([]byte("tuples-4"))
					v = b4.Get([]byte(string(partial)))
				} else if string(partial[0]) <= "o" {
					b5 := tx.Bucket([]byte("tuples-5"))
					v = b5.Get([]byte(string(partial)))
				} else if string(partial[0]) <= "r" {
					b6 := tx.Bucket([]byte("tuples-6"))
					v = b6.Get([]byte(string(partial)))
				} else if string(partial[0]) <= "u" {
					b7 := tx.Bucket([]byte("tuples-7"))
					v = b7.Get([]byte(string(partial)))
				} else {
					b8 := tx.Bucket([]byte("tuples-8"))
					v = b8.Get([]byte(string(partial)))
				}

				vals := string(v)
				// log.Println(partial)
				// log.Printf("The answer is: %v\n", vals)
				if len(v) > 0 {
					for _, k := range strings.Split(vals, " ") {
						db.View(func(tx *bolt.Tx) error {
							var v1 []byte
							if string(k[0]) <= "c" { // DIVIDED 6x: 32MB 84 ms...UNDIVIDED: 188ms
								b1 := tx.Bucket([]byte("words-1"))
								v1 = b1.Get([]byte(k))
							} else if string(k[0]) <= "f" {
								b2 := tx.Bucket([]byte("words-2"))
								v1 = b2.Get([]byte(k))
							} else if string(k[0]) <= "i" {
								b3 := tx.Bucket([]byte("words-3"))
								v1 = b3.Get([]byte(k))
							} else if string(k[0]) <= "l" {
								b4 := tx.Bucket([]byte("words-4"))
								v1 = b4.Get([]byte(k))
							} else if string(k[0]) <= "o" {
								b5 := tx.Bucket([]byte("words-5"))
								v1 = b5.Get([]byte(k))
							} else if string(k[0]) <= "r" {
								b6 := tx.Bucket([]byte("words-6"))
								v1 = b6.Get([]byte(k))
							} else if string(k[0]) <= "u" {
								b7 := tx.Bucket([]byte("words-7"))
								v1 = b7.Get([]byte(k))
							} else {
								b8 := tx.Bucket([]byte("words-8"))
								v1 = b8.Get([]byte(k))
							}
							v := string(v1)

							// v := string(b.Get([]byte(k)))
							// b := tx.Bucket([]byte("words"))
							_, ok := matches[v]
							if ok != true {
								matches[v] = levenshtein.Distance(s, v)
								// fmt.Printf("Word match: %v\n", v)
								// fmt.Printf("Distance : %v\n", matches[v])
							}
							return nil
						})
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

	return bestMatch, bestVal
}
