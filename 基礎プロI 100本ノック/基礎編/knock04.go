// No. 04 入力と計算
// 整数値を入力させ、その入力値を3倍した計算結果を表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	fmt.Print("input number: ")
	var x int
	fmt.Scan(&x)

	fmt.Println("answer =", x*3)
}
