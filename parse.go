package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cheggaaa/pb"
)

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
	fmt.Println("Parsing subsets...")
	bar := pb.StartNew(totalLines)
	for scanner.Scan() {
		bar.Increment()
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
	bar.FinishPrint("Finished parsing subsets")
	return
}

func dumpToBoltDB(path string, words map[string]int, tuples map[string]string, tupleLength int) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err == nil {
			os.Remove(path)
		}
	}

	// Open a new bolt database
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples-1"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples-2"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples-3"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples-4"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples-5"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples-6"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples-7"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("tuples-8"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("words"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
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
	start := time.Now()
	bar2 := pb.StartNew(len(words))
	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("words"))
		for k, v := range words {
			bar2.Increment()
			if len(k) > 0 {
				b.Put([]byte(strconv.Itoa(v)), []byte(k))
			}
		}
		return nil
	})
	elapsed := time.Since(start)
	bar2.FinishPrint("Words took " + elapsed.String())

	// fmt.Printf("inserting into bucket (tuple,words) '%v':'%v');\n", k, v)
	fmt.Println("Loading subsets into db...")
	start = time.Now()
	bar1 := pb.StartNew(len(tuples))
	bNums := []int{0, 0, 0, 0, 0, 0, 0, 0}
	db.Batch(func(tx *bolt.Tx) error {
		b1 := tx.Bucket([]byte("tuples-1"))
		b2 := tx.Bucket([]byte("tuples-2"))
		b3 := tx.Bucket([]byte("tuples-3"))
		b4 := tx.Bucket([]byte("tuples-4"))
		b5 := tx.Bucket([]byte("tuples-5"))
		b6 := tx.Bucket([]byte("tuples-6"))
		b7 := tx.Bucket([]byte("tuples-7"))
		b8 := tx.Bucket([]byte("tuples-8"))
		for k, v := range tuples {
			bar1.Increment()
			if string(k[0]) <= "c" { // DIVIDED 6x: 32MB 84 ms...UNDIVIDED: 188ms
				b1.Put([]byte(k), []byte(v))
				bNums[0]++
			} else if string(k[0]) <= "f" {
				b2.Put([]byte(k), []byte(v))
				bNums[1]++
			} else if string(k[0]) <= "i" {
				b3.Put([]byte(k), []byte(v))
				bNums[2]++
			} else if string(k[0]) <= "l" {
				b4.Put([]byte(k), []byte(v))
				bNums[3]++
			} else if string(k[0]) <= "o" {
				b5.Put([]byte(k), []byte(v))
				bNums[4]++
			} else if string(k[0]) <= "r" {
				b6.Put([]byte(k), []byte(v))
				bNums[5]++
			} else if string(k[0]) <= "u" {
				b7.Put([]byte(k), []byte(v))
				bNums[6]++
			} else {
				b8.Put([]byte(k), []byte(v))
				bNums[7]++
			}
		}
		return nil
	})
	elapsed = time.Since(start)
	bar1.FinishPrint("Subsets took " + elapsed.String())

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("vars"))
		err := b.Put([]byte("tupleLength"), []byte(strconv.Itoa(tupleLength)))
		return err
	})
}
func generateDB(wordpath string, path string, tupleLength int) {

	words, tuples := scanWords(wordpath, path, tupleLength)
	dumpToBoltDB(path, words, tuples, tupleLength)

}
