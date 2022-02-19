// No. 29 5つの合計
// 整数値を5回入力させ、それらの値の合計を表示するプログラムを繰り返しを使って作成せよ。
package main

import "fmt"

func main() {
	sum := 0
	var n int
	for i := 0; i < 5; i++ {
		fmt.Print("input number: ")
		fmt.Scan(&n)
		sum += n
	}
	fmt.Println("sum =", sum)
}
