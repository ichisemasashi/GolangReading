// No. 07 0 or not 0
// 整数値を入力させ、値が0ならzero、0でなければnot zeroと表示するプログラムを作成せよ。

package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if n == 0 {
		fmt.Println("zero")
	} else {
		fmt.Println("not zero")
	}
}
