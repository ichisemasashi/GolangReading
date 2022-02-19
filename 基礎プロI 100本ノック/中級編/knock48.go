// No. 48 繰り返し計算
// 最初に2以上の整数値を入力し、次の規則で計算し、計算回数と計算結果を表示し、計算結果が1になるまで繰り返すプログラムを作成せよ。
// 規則：ある値が偶数ならその値を1/2にする。奇数ならその値を3倍して1を足す。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("input number: ")
	fmt.Scan(&n)
	calc(n)
}

func calc(n int) {
	for i := 1; ; i++ {
		fmt.Printf("%v: ", i)
		if (n % 2) == 0 {
			n = n / 2
		} else {
			n = n*3 + 1
		}
		fmt.Println(n)
		if n == 1 {
			break
		}
	}
}
