// No. 33 入力値抜き
// 整数値を入力させ、1から9まで、入力値以外を表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	for i := 1; i <= 9; i++ {
		if i != n {
			fmt.Println(i)
		}
	}
}
