package goagrep

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cheggaaa/pb"
)

// VERBOSE is a flag to turn on/off status information during parsing
var VERBOSE bool

func init() {
	VERBOSE = false
}

func getPartials(s string, tupleLength int) []string {
	partials := make([]string, 100)
	num := 0
	s = strings.TrimSpace(strings.Replace(strings.ToLower(s), " ", "", -1))
	slen := len(s)
	if slen <= tupleLength {
		if slen <= 3 {
			partials[num] = "zzz"
			num = num + 1
		} else {
			for i := 0; i <= slen-3; i++ {
				partials[num] = s[i : i+3]
				num = num + 1
			}
		}
	} else {
		for i := 0; i <= slen-tupleLength; i++ {
			partials[num] = s[i : i+tupleLength]
			num = num + 1
			if num >= 100 {
				break
			}
		}
	}
	return partials[0:num]
}

func scanWords(wordpath string, tupleLength int, makeLookup bool) (words map[string]int, tuples map[string]string, wordsLookup map[int]string, tuplesLookup map[string][]int) {
	totalLines := lineCount(wordpath)

	inFile, _ := os.Open(wordpath)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	// initialize
	words = make(map[string]int)
	tuples = make(map[string]string)
	wordsLookup = make(map[int]string)
	tuplesLookup = make(map[string][]int)

	numTuples := 0
	numWords := 0
	lineNum := 0
	var bar *pb.ProgressBar
	if VERBOSE {
		fmt.Println("Parsing subsets...")
		bar = pb.StartNew(totalLines)
	}
	for scanner.Scan() {
		if VERBOSE {
			bar.Increment()
		}
		lineNum++
		s := strings.TrimSpace(scanner.Text())

		_, ok := words[s]
		if ok == false {
			if makeLookup {
				wordsLookup[numWords] = s
			} else {
				words[s] = numWords
			}

			partials := getPartials(s, tupleLength)
			for i := 0; i < len(partials); i++ {
				_, ok := tuples[partials[i]]
				if makeLookup {
					_, ok = tuplesLookup[partials[i]]
				}
				if ok == false {
					if makeLookup {
						tuplesLookup[partials[i]] = append([]int{}, numWords)
					} else {
						tuples[partials[i]] = strconv.Itoa(numWords)
					}
					numTuples++
				} else {
					if makeLookup {
						tuplesLookup[partials[i]] = append(tuplesLookup[partials[i]], numWords)
					} else {
						tuples[partials[i]] += " " + strconv.Itoa(numWords)
					}
				}
			}

			numWords++
		}

	}
	if VERBOSE {
		bar.FinishPrint("Finished parsing subsets")
	}
	return
}

func dumpToBoltDB(path string, words map[string]int, tuples map[string]string, tupleLength int) {
	var bar *pb.ProgressBar
	var start time.Time
	wordBuckets := int(len(words) / 600)
	if wordBuckets < 10 {
		wordBuckets = 10
	}
	if VERBOSE {
		fmt.Printf("Creating %v word buckets\n", wordBuckets)
	}

	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
		if VERBOSE {
			fmt.Println("Removed old " + path)
		}
	}

	// Open a new bolt database
	db, err := bolt.Open(path, 0600, &bolt.Options{NoGrowSync: false})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if VERBOSE {
		fmt.Println("Creating subset buckets...")
		bar = pb.StartNew(len(tuples))
		start = time.Now()
	}
	err = db.Batch(func(tx *bolt.Tx) error {
		for k := range tuples {
			if VERBOSE {
				bar.Increment()
			}
			firstLetter := string(k[0])
			secondLetter := string(k[1])
			if strings.Contains(alphabet, firstLetter) && strings.Contains(alphabet, secondLetter) {
				_, err := tx.CreateBucketIfNotExists([]byte("tuples-" + firstLetter + secondLetter))
				if err != nil {
					return fmt.Errorf("create bucket: %s", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if VERBOSE {
		elapsed := time.Since(start)
		bar.FinishPrint("Creating subset buckets took " + elapsed.String())
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if VERBOSE {
		fmt.Println("Creating words buckets...")
	}
	db.Batch(func(tx *bolt.Tx) error {
		for i := 0; i < wordBuckets; i++ {
			_, err := tx.CreateBucket([]byte("words-" + strconv.Itoa(i)))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("vars"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	// fmt.Printf("INSERT INTO words (id,word) values (%v,'%v');\n", v, k)
	if VERBOSE {
		fmt.Println("Loading words into db...")
		start = time.Now()
		bar = pb.StartNew(len(words))
	}
	err = db.Batch(func(tx *bolt.Tx) error {
		for k, v := range words {
			if VERBOSE {
				bar.Increment()
			}
			if len(k) > 0 {
				b := tx.Bucket([]byte("words-" + strconv.Itoa(int(math.Mod(float64(v), float64(wordBuckets))))))
				b.Put([]byte(strconv.Itoa(v)), []byte(k))
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if VERBOSE {
		elapsed := time.Since(start)
		bar.FinishPrint("Words took " + elapsed.String())
	}

	if VERBOSE {
		fmt.Println("Loading subsets into db...")
		start = time.Now()
		bar = pb.StartNew(len(tuples))
	}
	err = db.Update(func(tx *bolt.Tx) error {
		for k, v := range tuples {
			if VERBOSE {
				bar.Increment()
			}
			firstLetter := string(k[0])
			secondLetter := string(k[1])
			if strings.Contains(alphabet, firstLetter) && strings.Contains(alphabet, secondLetter) {
				b := tx.Bucket([]byte("tuples-" + firstLetter + secondLetter))
				b.Put([]byte(k), []byte(v))
			} else {
				b := tx.Bucket([]byte("tuples"))
				b.Put([]byte(k), []byte(v))
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err) // BUG(schollz): Windows file resize error: https://github.com/schollz/goagrep/issues/6
	}
	if VERBOSE {
		elapsed := time.Since(start)
		bar.FinishPrint("Subsets took " + elapsed.String())
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("vars"))
		err := b.Put([]byte("tupleLength"), []byte(strconv.Itoa(tupleLength)))
		return err
	})

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("vars"))
		err := b.Put([]byte("wordBuckets"), []byte(strconv.Itoa(wordBuckets)))
		return err
	})
}

// GenerateDB generates the database with precomputed strings for later searching.
// It is required to you generate a database before you use the GetMatch() function.
//
// stringListPath is the filename of the list of strings you want to use
//
// databasePath is the filename of the database that is outputed
//
// tupleLength is the length of the subsets you want to use
func GenerateDB(stringListPath string, databasePath string, tupleLength int, verbosity bool) {
	VERBOSE = verbosity
	words, tuples, _, _ := scanWords(stringListPath, tupleLength, false)
	dumpToBoltDB(databasePath, words, tuples, tupleLength)
}

func GenerateDBInMemory(stringListPath string, tupleLength int, verbosity bool) (words map[int]string, tuples map[string][]int) {
	VERBOSE = verbosity
	_, _, words, tuples = scanWords(stringListPath, tupleLength, true)
	return
}
