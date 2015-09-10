package main

import (
    "fmt"
    "github.com/arbovm/levenshtein"
    "bufio"
    "os"
)


func addToFile(s string) {
    slen := len(s)
    if slen <= 4 {
       addToCache("four",s)
    } else {
        for i := 0; i <= slen-3; i ++ {
            addToCache(s[i:i+3],s)
        }
    }
}

func addToCache(spartial string, s string) {
    f, err := os.OpenFile("cache/" + spartial, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    if _, err = f.WriteString(s); err != nil {
        panic(err)
    }
}


func generateHash(path string) {
  inFile, _ := os.Open(path)
  defer inFile.Close()
  scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines) 
  
  for scanner.Scan() {
    addToFile(scanner.Text())
  }
}

func main() {
    generateHash("wordlist")
    s1 := "marcy playground"
    s2 := "mary playground"
    fmt.Printf("The distance between %v and %v is %v\n",
        s1, s2, levenshtein.Distance(s1, s2))
    // -> The distance between kitten and sitting is 3
}
