// No. 43 2次方程式の解の判別
// 2次方程式 ax^2 + bx + c = 0 （x^2はxの2乗の意味）の係数a, b, cを入力し、2次方程式の解が2つの実数解か、重解か、2つの虚数解かを判別するプログラムを作成せよ。
package main

import "fmt"

func main() {
	a := readNumber("input a: ")
	b := readNumber("input b: ")
	c := readNumber("input c: ")

	d := b*b - 4*a*c
	if d < 0 {
		fmt.Println("2つの虚数解")
	} else if d == 0 {
		fmt.Println("重解")
	} else {
		fmt.Println("2つの実数解")
	}
}

func readNumber(p string) int {
	var n int
	fmt.Print(p)
	fmt.Scan(&n)
	return n
}
