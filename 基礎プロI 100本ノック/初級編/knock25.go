// No. 25 -10未満?、-10以上0未満?、0以上?
// 整数値を入力させ、その値が-10未満ならrange 1、-10以上0未満であればrange 2、0以上であればrange 3、と表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if n < -10 {
		fmt.Println("range 1")
	} else if (-10 <= n) && (n < 0) {
		fmt.Println("range 2")
	} else {
		fmt.Println("range 3")
	}
}
