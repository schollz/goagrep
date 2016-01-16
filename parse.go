package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

func getPartials(s string, tupleLength int) []string {
	partials := make([]string, 100000)
	num := 0
	s = strings.ToLower(s)
	s = strings.Replace(s, "/", "", -1)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "the", "", -1)
	s = strings.Replace(s, "by", "", -1)
	s = strings.Replace(s, "dr", "", -1)
	s = strings.Replace(s, "of", "", -1)
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

func scanWords(wordpath string, path string, tupleLength int) (words map[string]int, tuples map[string]string) {

	inFile, _ := os.Open(wordpath)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	words = make(map[string]int)
	tuples = make(map[string]string)
	numTuples := 0
	numWords := 0

	lineNum := 0

	for scanner.Scan() {
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
		_, err := tx.CreateBucket([]byte("tuples"))
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
	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("words"))
		for k, v := range words {
			b.Put([]byte(strconv.Itoa(v)), []byte(k))
		}
		return nil
	})

	// fmt.Printf("inserting into bucket uples '%v':'%v');\n", k, v)
	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tuples"))
		for k, v := range tuples {
			b.Put([]byte(k), []byte(v))
		}
		return nil
	})

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
