// No. 16 0でおしまい
// 整数値を入力させ、入力値が0でなければ再度入力させ、0であれば終了するプログラムを作成せよ。
package main

import "fmt"

func readNumber(prom string) int {
	var n int
	fmt.Print(prom)
	fmt.Scan(&n)
	return n
}

func main() {
	n := readNumber("input number: ")
	for n != 0 {
		n = readNumber("input number: ")
	}
}
