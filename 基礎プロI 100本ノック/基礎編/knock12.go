// No. 12 ごあいさつ指定回
// 整数値を入力させ、その値の回数だけHello World!を繰り返して表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		fmt.Println("Hello World!")
	}
}
