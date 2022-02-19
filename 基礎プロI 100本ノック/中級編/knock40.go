// No. 40 even or odd
// 整数値を入力させ、その値が偶数ならばeven、奇数ならばoddと表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if (n % 2) == 0 {
		fmt.Printf("%v is even.\n", n)
	} else {
		fmt.Printf("%v is odd.\n", n)
	}
}
