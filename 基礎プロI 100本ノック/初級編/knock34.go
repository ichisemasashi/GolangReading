// No. 34 入力値抜き改
// 整数値を入力させ、1から9まで、入力値と入力値+1以外を表示するプログラムを作成せよ。入力値が9の場合は9のみ表示しない。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	for i := 1; i <= 9; i++ {
		if (i == n) || (i == (n + 1)) {

		} else {
			fmt.Println(i)
		}
	}
}
