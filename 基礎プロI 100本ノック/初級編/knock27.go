// No. 27 1からnまでの和
// 整数値を入力させ、1からその値までの総和を計算して表示するプログラムを作成せよ。ただし、0以下の値を入力した場合は0と表示する。
package main

import (
	"fmt"
)

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	sum := 0
	for i := 0; i <= n; i++ {
		sum += i
	}

	fmt.Println("sum =", sum)
}
