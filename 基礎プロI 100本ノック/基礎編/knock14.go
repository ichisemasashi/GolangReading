// No. 14 カウントダウン
// 整数値を入力させ、入力値から0まで数を1ずつ減らして表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	for ; 0 <= n; n-- {
		fmt.Println(n)
	}
}
