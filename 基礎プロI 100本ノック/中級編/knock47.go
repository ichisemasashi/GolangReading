// No. 47 値の入れ替え
// 異なる整数値を2つ入力し、それぞれ別の変数に格納する。そして、それらの変数の値を入れ替えた後、それぞれの変数の値を表示するプログラムを作成せよ。単に順序を変えて表示するだけではダメ。必ず変数の値を入れ替えること。下の実行例の場合、まず変数aに2、bに5が入力された状態から、aの値が5、bの値が2になるように入れ替える。
package main

import "fmt"

func main() {
	a := readNumber("input a: ")
	b := readNumber("input b: ")

	a, b = b, a
	fmt.Printf("a = %v, b = %v\n", a, b)
}

func readNumber(p string) int {
	var n int
	fmt.Print(p)
	fmt.Scan(&n)
	return n
}
