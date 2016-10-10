package goagrep

import (
	"fmt"
	"strings"
)

func ExampleUtils1() {
	fmt.Println(removeDuplicates([]int{1, 2, 3, 4, 5, 5, 5, 6, 7, 7, 8, 10}))
	// Output: [1 2 3 4 5 6 7 8 10]
}

func ExampleUtils2() {
	fmt.Println(abs(-1), abs(30))
	// Output: 1 30
}

func ExampleUtils3() {
	fmt.Println(stringInSlice("hello", []string{"hello", "hell", "hello there"}), stringInSlice("hello", []string{"hell", "hello there"}))
	// Output: true false
}

func ExampleUtils4() {
	fmt.Println(lineCount("../example/testlist"))
	// Output: 1009
}

func ExampleUtils5() {
	fmt.Println(getDistance("bread", "bed"))
	// Output: 2
}

func ExampleUtils6() {
	Normalize = true
	fmt.Println(getDistance("Italy Luxury: smething something", "Italy Luxury"))
	Normalize = false
	// Output: 0
}

func ExampleUtils7() {
	s1, subsets := getSubstrings("harry", "a great big world")
	fmt.Println(s1, strings.Join(subsets, ","))
	// Output: harry a gre, grea,great,reat ,eat b,at bi,t big, big ,big w,ig wo,g wor, worl,world
}
