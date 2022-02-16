// No. 06 0?
// 整数値を入力させ、値が0ならzeroと表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if n == 0 {
		fmt.Println("zero")
	}
}
