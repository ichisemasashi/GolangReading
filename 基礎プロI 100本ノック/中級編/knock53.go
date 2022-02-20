// No. 53 素因数分解
// 自然数の入力値を素因数分解して結果を表示するプログラムを作成せよ。
// * まず2で割り切れる間は割っていき、2で割り切れなくなったら3で、と、割る数を1ずつ大きくしながら繰り返す。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)

	for i := 2; n > 1; i++ {
		for {
			m := n % i
			if m == 0 {
				fmt.Print(i, " ")
			} else {
				break
			}
			n = int(n / i)
		}
	}
	fmt.Println("")
}
