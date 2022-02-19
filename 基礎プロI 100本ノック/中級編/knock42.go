// No. 42 小さい順?
// 整数値を3つ入力させ、それらの値が小さい順（等しくてもよい）に並んでいるか判定し、小さい順に並んでいる場合はOK、そうなっていない場合はNGと表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	n1 := readNumber("input number 1: ")
	n2 := readNumber("input number 2: ")
	n3 := readNumber("input number 3: ")

	if (n1 <= n2) && (n2 <= n3) {
		fmt.Println("OK")
	} else {
		fmt.Println("NG")
	}
}

func readNumber(p string) int {
	var n int
	fmt.Print(p)
	fmt.Scan(&n)
	return n
}
