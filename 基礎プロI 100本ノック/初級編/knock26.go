// No. 26 switch-case
// 整数値を入力させ、その値が1ならone、2ならtwo、3ならthree、それ以外ならothersと表示するプログラムをswicth-case文を使って作成せよ。
package main

import (
	"fmt"
)

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	var s string
	switch n {
	case 1:
		s = "one"
	case 2:
		s = "two"
	case 3:
		s = "three"
	default:
		s = "other"
	}
	fmt.Println(s)
}
