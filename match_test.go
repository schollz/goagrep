package main

import "testing"

func BenchmarkMatch(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getMatch2("heroint", "testlist.db")
	}
}

//  go build match.go &&  ./match build testlist | sqlite3 testlist.db && go test -bench=.
