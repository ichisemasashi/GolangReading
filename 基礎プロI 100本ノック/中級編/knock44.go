// No. 44 通貨換算
// 換算したい金額（円単位）と1ドル何円かを整数値で入力すると、入力した金額が何ドル何セントか計算して表示するプログラムを作成せよ。1セントは1/100ドルである。結果は整数値でよい（1セント未満の端数切り捨て）。
package main

import "fmt"

func main() {
	x := readNumber("何円? ")
	y := readNumber("1ドルは何円? ")
	daller := int(x / y)
	cent := int(x*100/y - daller*100)
	fmt.Printf("%v円は%vドル%vセント\n", x, daller, cent)
}

func readNumber(p string) int {
	var n int
	fmt.Print(p)
	fmt.Scan(&n)
	return n
}
