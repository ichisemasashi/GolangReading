// No. 20 割って掛ける
// 整数値を2つ入力させ、1つめの値を2つめの値で割った結果を表示し、続けてその結果に2つめの値を掛けた結果を表示するプログラムを作成せよ。計算はすべて整数型で行うこと（割り切れない場合は自動的に小数点以下が切り捨てられる）。
package main

import "fmt"

func main() {
	x := readNumber("input 1st value: ")
	y := readNumber("input 2nd value: ")

	// 割り算
	r := int(x / y)
	fmt.Println("result:", r)
	// かけ算
	r = r * y
	fmt.Println("result:", r)
}

func readNumber(p string) int {
	var n int
	fmt.Print(p)
	fmt.Scan(&n)
	return n
}
