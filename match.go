package main

import (
	"github.com/arbovm/levenshtein"
	"github.com/cheggaaa/pb"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

//GLOBALS
var findings_matches []string
var findings_leven []int
var wg sync.WaitGroup
var tuple_length int
var file_tuple_length int

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

func generateHash(path string) {
	inFile2, _ := os.Open(path)
	numLines, _ := lineCounter(inFile2)
	inFile2.Close()

	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	mm := make(map[string]string)

	fmt.Printf("Building map...\n")
	bar := pb.StartNew(numLines)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		bar.Increment()
		s := strings.Replace(scanner.Text(), "/", "", -1)
		//addToCache("keys.list", s)
		partials := getPartials(s)
		for i := 0; i < len(partials); i++ {
			_, ok := mm[partials[i]]
			if ok == true {
				mm[partials[i]] = mm[partials[i]] + " " + strconv.Itoa(lineNum)
			} else {
				mm[partials[i]] = strconv.Itoa(lineNum)

			}
			//addToCache(partials[i], strconv.Itoa(lineNum))
		}
	}
	bar.FinishPrint("Finished.\n")
	fmt.Printf("Building cache...")
	for k := range mm {
		//fmt.Printf("%v : %v\n", k, mm[k])
		addToCache(k[0:file_tuple_length], k+" "+mm[k])
	}
}

func addToCache(spartial string, s string) {
	f, err := os.OpenFile("cache/"+spartial, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(s + "\n"); err != nil {
		panic(err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
			break
		}
	}
	return false
}

func getPartials(s string) []string {
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
	if slen <= tuple_length {
		partials[num] = "zzzf"
		num = num + 1
	} else {
		for i := 0; i <= slen-tuple_length; i++ {
			partials[num] = s[i : i+tuple_length]
			num = num + 1
		}
	}
	return partials[0:num]
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

func ReadLine(file string, lineNum int) (line string, lastLine int, err error) {
	r, _ := os.Open(file)
	defer r.Close()
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

func getIndiciesFromPartial(partials []string, path string) []int {
	indexMatches := make([]int, 100000)
	numm := 0
	for i := 0; i < len(partials); i++ {

		inFile, _ := os.Open(path + partials[i][0:file_tuple_length])
		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			scan := scanner.Text()
			if partials[i] == scan[0:tuple_length] {
				for _, k := range strings.Split(scan[tuple_length:], " ") {
					indexMatches[numm], _ = strconv.Atoi(k)
					numm++
				}
			}
		}

	}
	//fmt.Printf("\nIndex matches: %v\n", indexMatches[0:numm])
	indexMatches = removeDuplicates(indexMatches[0:numm])
	//fmt.Printf("\nIndex matches: %v\n", indexMatches)
	return indexMatches

}

func getMatch(s string, path string) (string, int) {
	start := time.Now()
	partials := getPartials(s)
	elapsed := time.Since(start)
	fmt.Printf("Partials took %s", elapsed)
	//fmt.Printf("Partials: %v", partials)
	runtime.GOMAXPROCS(8)
	N := 8

	start = time.Now()
	indexMatches := getIndiciesFromPartial(partials, path)
	fmt.Printf("Indices from partials took %s", time.Since(start))
	
	matches := make([]string, len(indexMatches))
	for i := 0; i < len(indexMatches); i++ {
		if indexMatches[i] > 0 {
			str, _, err := ReadLine("cache/keys.list", indexMatches[i])
			//fmt.Printf("\n--%v %v %v--\n", str, i, indexMatches[i])
			if err != nil {
				//fmt.Printf("Error reading line ", lastLine)
			}
			matches[i] = str
		}
	}

	//fmt.Printf("%v\n", matches)

	findings_leven = make([]int, N)
	findings_matches = make([]string, N)

	wg.Add(N)
	for i := 0; i < N; i++ {
		go search(matches[i*len(matches)/N:(i+1)*len(matches)/N], s, i)
	}
	wg.Wait()

	lowest := 100
	best_index := 0
	for i := 0; i < len(findings_leven); i++ {
		if findings_leven[i] < lowest {
			lowest = findings_leven[i]
			best_index = i
		}
	}

	return findings_matches[best_index], lowest
}

func search(matches []string, target string, process int) {
	defer wg.Done()
	match := "No match"
	target = strings.ToLower(target)
	bestLevenshtein := 1000
	for i := 0; i < len(matches); i++ {
		d := levenshtein.Distance(target, strings.ToLower(matches[i]))
		if d < bestLevenshtein {
			bestLevenshtein = d
			match = matches[i]
		}
	}
	findings_matches[process] = match
	findings_leven[process] = bestLevenshtein
}

func main() {
	tuple_length = 6
	file_tuple_length = 4
	if strings.EqualFold(os.Args[1], "help") {
		fmt.Printf("Version 1.1 - %v-mer tuples, removing commons\n", tuple_length)
		fmt.Println("./match-concurrent build <NAME OF WORDLIST> - builds cache/ folder in current directory\n")
		fmt.Println("./match-concurrent 'word or words to match' /directions/to/cache/\n")
	} else if strings.EqualFold(os.Args[1], "build") {
		os.Mkdir("cache", 0775)
		generateHash(os.Args[2])
		// open files r and w
		r, err := os.Open(os.Args[2])
		if err != nil {
			panic(err)
		}
		defer r.Close()

		w, err := os.Create("cache/keys.list")
		if err != nil {
			panic(err)
		}
		defer w.Close()

		// do the actual work
		n, err := io.Copy(w, r)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Copied %v bytes\n", n)

	} else {
		match, lowest := getMatch(os.Args[1], os.Args[2])
		fmt.Printf("%v|||%v\n", match, lowest)
	}
}
