// No. 18 配列を入力値で初期化
// 要素数10の整数型の配列を宣言し、整数値を入力させ、すべての配列の要素を入力値として、すべての要素の値を表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	a := [10]int{}
	for i := 0; i < 10; i++ {
		a[i] = n
	}

	for _, x := range a {
		fmt.Println(x)
	}
}
