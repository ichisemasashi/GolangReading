// No. 28 nの階乗
// 整数値を入力させ、その値の階乗を表示するプログラムを作成せよ。ただし、0以下の値を入力した場合は1と表示する。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	var p uint64
	if n <= 0 {
		p = 0
	} else {
		p = fact(n)
	}
	fmt.Println("factorial =", p)
}

func fact(n int) uint64 {
	res := uint64(1)
	for i := 1; i < n; i++ {
		res *= uint64(i)
	}
	return res
}
