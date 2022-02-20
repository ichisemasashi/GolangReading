// No. 56 2進数変換
// 0〜65535の整数値を入力させ、入力値を16桁の2進数に変換して表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n uint
	fmt.Print("input number: ")
	fmt.Scan(&n)
	fmt.Printf("%016b\n", n)
}
