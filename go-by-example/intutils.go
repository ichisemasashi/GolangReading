package main

import "fmt"

// この整数の最小値を求めるシンプルな実装をテストしたい
// としましょう。 通常、テストしたいコードは intutils.go
// のような名前の ソースファイル中にあり、そのテストファイル
// は intutils_test.go のような名前になります。
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	i := IntMin(1, 2)
	fmt.Println(i)
}
