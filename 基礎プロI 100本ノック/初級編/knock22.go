// No. 22 -10以下または10以上
// 整数値を入力させ、その値が-10以下、または、10以上であればOKと表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if (n <= -10) || (10 <= n) {
		fmt.Println("OK")
	}
}
