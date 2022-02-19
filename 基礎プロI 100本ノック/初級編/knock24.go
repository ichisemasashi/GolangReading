// No. 24 -10以上0未満、または、10以上
// 整数値を入力させ、その値が-10以上0未満、または、10以上であればOK、そうでなければNGと表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if ((-10 < n) && (n < 0)) || (10 <= n) {
		fmt.Println("OK")
	} else {
		fmt.Println("NG")
	}
}
