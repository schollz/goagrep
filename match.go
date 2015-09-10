package main

import (
    "fmt"
    "github.com/arbovm/levenshtein"
    "bufio"
    "os"
)


func addToHash(m map[string][]string, s string) {
    slen := len(s)
    if slen <= 4 {
        m["four"] = append(m["four"], s)

    } else {

    for i := 0; i <= slen-3; i ++ {
        fmt.Println(s[i:i+3])
        m[s[i:i+3]] = append(m[s[i:i+3]], s)
    }

    }
}


func generateHash(path string, m map[string][]string) {
  inFile, _ := os.Open(path)
  defer inFile.Close()
  scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines) 
  
  for scanner.Scan() {
    addToHash(m,scanner.Text())
  }
}

func main() {
    m := make(map[string][]string)
    generateHash("wordlist",m)
    addToHash(m,"something")
    addToHash(m,"some")
for key, value := range m {
    fmt.Println("Key:", key, "Value:", value)
}
    s1 := "marcy playground"
    s2 := "mary playground"
    fmt.Printf("The distance between %v and %v is %v\n",
        s1, s2, levenshtein.Distance(s1, s2))
    // -> The distance between kitten and sitting is 3
}
