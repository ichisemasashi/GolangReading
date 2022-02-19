// No. 30 棒グラフ
// 整数値を入力させ、その個数だけ*を表示するプログラムを作成せよ。入力値が0以下の値の場合は何も書かなくてよい。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	printStars(n)
}

func printStars(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("*")
	}
	if 0 < n {
		fmt.Println("")
	}

}
