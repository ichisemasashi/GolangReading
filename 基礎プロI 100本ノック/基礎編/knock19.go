// No. 19 配列に入力値を格納
// 要素数5の整数型の配列を宣言し、すべての配列に対して順に入力された整数値を代入し、すべての要素の値を表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	a := [5]int{}
	for i, _ := range a {
		fmt.Print("input number: ")
		fmt.Scan(&(a[i]))
	}

	for _, x := range a {
		fmt.Println(x)
	}
}
