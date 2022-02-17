// No. 10 絶対値
// 整数値を入力させ、その値を絶対値にして表示するプログラムを作成せよ。（できれば変数の値を絶対値に変えるようにせよ）

package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if 0 <= n {
		fmt.Println("absolute value is", n)
	} else {
		fmt.Println("absolute value is", -n)
	}
}
