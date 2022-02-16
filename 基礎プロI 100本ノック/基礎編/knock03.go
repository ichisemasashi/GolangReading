// No. 03 入力
// 整数値を入力させ、その入力値を表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	fmt.Print("input number: ")
	var x int
	fmt.Scan(&x)

	fmt.Println("your number is", x)
}
