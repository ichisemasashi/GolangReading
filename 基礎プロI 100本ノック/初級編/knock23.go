// No. 23 -5以上10未満
// 整数値を入力させ、その値が-5以上10未満であればOK、そうでなければNGと表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if (-5 <= n) && (n < 10) {
		fmt.Println("OK")
	} else {
		fmt.Println("NG")
	}
}
