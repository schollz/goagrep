package goagrep

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	"github.com/arbovm/levenshtein"
)

var alphabet string

func init() {
	alphabet = "abcdefghijklmnopqrstuvwxyz-"
}

func getDistance(s1 string, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)
	// dist := hamming.Calc(s1, s2)
	// if dist > 0 {
	// 	return dist
	// } else {
	dist := levenshtein.Distance(s1, s2)
	if Normalize {
		s, substrings := getSubstrings(s1, s2)
		dist = 10000
		for _, sub := range substrings {
			distSub := levenshtein.Distance(s, sub)
			if distSub < dist {
				dist = distSub
			}
		}
	}
	return dist
}

func removeDuplicates(a []int) []int {
	result := []int{}
	seen := map[int]int{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

func getSubstrings(s1 string, s2 string) (string, []string) {
	check := ""
	tosplit := ""
	subsets := []string{}
	if len(s1) == len(s2) {
		return s1, append(subsets, s2)
	} else if len(s1) < len(s2) {
		check = s1
		tosplit = s2
	} else {
		check = s2
		tosplit = s1
	}
	i := 0
	splitSize := len(check)
	subsets = make([]string, len(tosplit)-len(check)+1)
	for {
		subsets[i] = tosplit[i : i+splitSize]
		i++
		if i+splitSize > len(tosplit) {
			break
		}
	}
	return check, subsets
}
