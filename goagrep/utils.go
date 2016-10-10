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
		minWord := len(s1)
		if len(s2) < minWord {
			minWord = len(s2)
		}
		if dist > minWord {
			dist = (dist-minWord)/2 + minWord
		}
		lcsDist := LCS(s1, s2)
		dist = dist + -1*lcsDist
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

func Max(more ...int) int {
	max_num := more[0]
	for _, elem := range more {
		if max_num < elem {
			max_num = elem
		}
	}
	return max_num
}

func LCS(str1, str2 string) int {
	len1 := len(str1)
	len2 := len(str2)

	table := make([][]int, len1+1)
	for i := range table {
		table[i] = make([]int, len2+1)
	}

	i, j := 0, 0
	for i = 0; i <= len1; i++ {
		for j = 0; j <= len2; j++ {
			if i == 0 || j == 0 {
				table[i][j] = 0
			} else if str1[i-1] == str2[j-1] {
				table[i][j] = table[i-1][j-1] + 1
			} else {
				table[i][j] = Max(table[i-1][j], table[i][j-1])
			}
		}
	}
	return table[len1][len2] //, Back(table, str1, str2, len1-1, len2-1)
}

//http://en.wikipedia.org/wiki/Longest_common_subsequence_problem
func Back(table [][]int, str1, str2 string, i, j int) string {
	if i == 0 || j == 0 {
		return ""
	} else if str1[i] == str2[j] {
		return Back(table, str1, str2, i-1, j-1) + string(str1[i])
	} else {
		if table[i][j-1] > table[i-1][j] {
			return Back(table, str1, str2, i, j-1)
		} else {
			return Back(table, str1, str2, i-1, j)
		}
	}
}
