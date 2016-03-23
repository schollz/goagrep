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

var VERBOSE bool

func init() {
	VERBOSE = true
}

func getPartials(s string, tupleLength int) []string {
	partials := make([]string, 500)
	num := 0
	s = strings.ToLower(s)
	s = strings.Replace(s, "/", "", -1)
	s = strings.Replace(s, " the ", "", -1)
	s = strings.Replace(s, " by ", "", -1)
	s = strings.Replace(s, " dr", "", -1)
	s = strings.Replace(s, " of ", "", -1)
	s = strings.Replace(s, " and ", "", -1)
	s = strings.Replace(s, " ", "", -1)
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
		}
	}
	//fmt.Println(s, partials[0:num])
	return partials[0:num]
}

func lineCount(filepath string) (numLines int) {
	// open a file
	numLines = 0
	if file, err := os.Open(filepath); err == nil {

		// make sure it gets closed
		defer file.Close()

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			numLines = numLines + 1
		}

		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal(err)
	}
	return
}

func scanWords(wordpath string, path string, tupleLength int) (words map[string]int, tuples map[string]string) {

	totalLines := lineCount(wordpath)

	inFile, _ := os.Open(wordpath)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	words = make(map[string]int)
	tuples = make(map[string]string)
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
		s := strings.Replace(scanner.Text(), "/", "", -1)
		s = strings.Replace(s, "'", "", -1)

		_, ok := words[s]
		if ok == false {
			words[s] = numWords

			partials := getPartials(s, tupleLength)
			for i := 0; i < len(partials); i++ {
				_, ok := tuples[partials[i]]
				if ok == false {
					tuples[partials[i]] = strconv.Itoa(numWords)
					numTuples++
				} else {
					tuples[partials[i]] += " " + strconv.Itoa(numWords)
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

	fmt.Println("Creating subset buckets...")
	bar3 := pb.StartNew(len(tuples))
	start := time.Now()
	err = db.Batch(func(tx *bolt.Tx) error {
		for k := range tuples {
			bar3.Increment()
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
	elapsed := time.Since(start)
	bar3.FinishPrint("Creating subset buckets took " + elapsed.String())

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	fmt.Println("Creating words buckets...")
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
	fmt.Println("Loading words into db...")
	start = time.Now()
	bar2 := pb.StartNew(len(words))
	err = db.Batch(func(tx *bolt.Tx) error {
		for k, v := range words {
			bar2.Increment()
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
	elapsed = time.Since(start)
	bar2.FinishPrint("Words took " + elapsed.String())

	// fmt.Printf("inserting into bucket (tuple,words) '%v':'%v');\n", k, v)
	fmt.Println("Loading subsets into db...")
	start = time.Now()
	bar1 := pb.StartNew(len(tuples))
	err = db.Update(func(tx *bolt.Tx) error {
		for k, v := range tuples {
			bar1.Increment()
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
	elapsed = time.Since(start)
	bar1.FinishPrint("Subsets took " + elapsed.String())

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
func GenerateDB(stringListPath string, databasePath string, tupleLength int) {

	words, tuples := scanWords(stringListPath, databasePath, tupleLength)
	dumpToBoltDB(databasePath, words, tuples, tupleLength)

}
