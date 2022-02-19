// No. 36 続・配列要素の参照
// {3, 7, 0, 8, 4, 1, 9, 6, 5, 2}で初期化される大きさ10の整数型配列を宣言し、整数値を2つ入力させ、要素番号が入力値である2つの配列要素の値の積を計算して表示するプログラムを作成せよ。入力値が配列の要素の範囲外であるかどうかのチェックは省略してよい。
package main

import "fmt"

func main() {
	x := readNumber("input index1: ")
	y := readNumber("input index2: ")

	a := [10]int{3, 7, 0, 8, 4, 1, 9, 6, 5, 2}
	fmt.Printf("%v * %v = %v\n", a[x], a[y], a[x]*a[y])
}

func readNumber(p string) int {
	var res int
	fmt.Print(p)
	fmt.Scan(&res)
	return res
}
