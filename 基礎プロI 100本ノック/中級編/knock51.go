// No. 51 お支払い
// 指定した金額を100円玉と10円玉と1円玉だけで、できるだけ少ない枚数で支払いたい。金額を入力するとそれぞれの枚数を計算して表示するプログラムを作成せよ。
// * 問題文は「できるだけ少ない」となっているが、金額の大きな硬貨を優先して使うように計算すれば自然とそうなる。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input money: ")
	fmt.Scan(&n)

	coins := []int{100, 10, 1}
	for _, c := range coins {
		m := int(n / c)
		fmt.Printf("%v円玉%v枚 ", c, m)
		n = n - m*c
	}

	fmt.Println("")
}
