// No. 13 カウントアップ
// 整数値を入力させ、0から入力値まで数を1ずつ増やして表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	for i := 0; i <= n; i++ {
		fmt.Println(i)
	}
}
