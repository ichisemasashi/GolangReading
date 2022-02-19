// No. 31 棒グラフ改
// 整数値を入力させ、その個数だけ*を、5個おきに空白（スペース）を入れて表示するプログラムを作成せよ。入力値が0以下の値の場合は何も書かなくてよい。
// ** No.30に加えて、カウンタ変数を5で割った余りが「特定の値」の場合にスペースを表示するようにすればよい。「特定の値」が何であるかは、繰り返し回数の数え方に依存する。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if 0 < n {
		printStars(n)
	}
}

func printStars(n int) {
	for i := 1; i <= n; i++ {
		if (i % 5) == 0 {
			fmt.Print("* ")
		} else {
			fmt.Print("*")
		}
	}
	fmt.Println("")
}
