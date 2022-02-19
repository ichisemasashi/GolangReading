// No. 32 5の倍数でbar
// 1から20まで順に表示するが、5の倍数の場合は数字の代わりにbarと表示するプログラムを作成せよ。
package main

import "fmt"

func main() {
	for i := 1; i <= 20; i++ {
		if (i % 5) == 0 {
			fmt.Println("bar")
		} else {
			fmt.Println(i)
		}
	}
}
