// No. 21 5より大きく20より小さい
// 整数値を入力させ、その値が5よりも大きく、かつ、20よりも小さければOKと表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if (5 < n) && (n < 20) {
		fmt.Println("OK")
	}
}
