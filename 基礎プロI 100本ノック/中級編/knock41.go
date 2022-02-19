// No. 41 1桁の自然数?
// 整数値を入力させ、その値が一桁の自然数かそうでないか判定するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if (0 < n) && (n < 10) {
		fmt.Printf("%v is a single figure.\n", n)
	} else {
		fmt.Printf("%v is not a single figure.\n", n)
	}
}
