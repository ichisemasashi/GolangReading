package main

import "fmt"

func main() {
	// 文字列。 + で互いに連結できます。
	fmt.Println("go" + "lang")

	// 整数や浮動小数。
	fmt.Println("1+1 =", 1+1)
	fmt.Println("7.0/3.0 =", 7.0/3.0)

	// 真偽値。ブール演算子も期待通り使えます。
	fmt.Println(true && true)
	fmt.Println(true || false)
	fmt.Println(!true)
}
