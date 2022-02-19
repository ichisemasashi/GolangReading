// No. 49 九九
// 九九の表を、2重の繰り返しを使って表示するプログラムを作成せよ。2重の繰り返しを使わず単に表示するだけではダメ。値の間はタブ(\t)を使って間をあけること。
package main

import "fmt"

func main() {
	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			fmt.Print(i*j, "\t")
		}
		fmt.Print("\n")
	}
}
