// No. 09 正? 負? 0?
// 整数値を入力させ、値が正であればpositive、負であればnegative、0であればzeroと表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	if n < 0 {
		fmt.Println("negative")
	} else if n == 0 {
		fmt.Println("zero")
	} else {
		fmt.Println("positive")
	}
}
