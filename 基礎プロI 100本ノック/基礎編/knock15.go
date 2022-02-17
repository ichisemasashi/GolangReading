// No. 15 2ずつカウントアップ
// 整数値を入力させ、0から入力値を超えない値まで2ずつ増やして表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	for i := 0; i < n; i = i + 2 {
		fmt.Println(i)
	}
}
