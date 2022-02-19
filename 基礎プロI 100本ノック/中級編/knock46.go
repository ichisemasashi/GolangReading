// No. 46 入場料
// 神山美術館の入場料金は、一人600円であるが、5人以上のグループなら一人550円、20人以上の団体なら一人500円である。人数を入力し、入場料の合計を計算するプログラムを作成せよ。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("人数: ")
	fmt.Scan(&n)

	pay := 0
	if n < 5 {
		pay = n * 600
	} else if (5 <= n) && (n < 20) {
		pay = n * 550
	} else {
		pay = n * 500
	}

	fmt.Println("料金", pay)
}
