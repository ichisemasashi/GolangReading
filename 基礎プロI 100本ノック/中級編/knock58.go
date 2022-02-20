// No. 58 棒グラフ
// 0以上の整数値を5つ入力させ、それぞれの値に対して、次の形式でグラフを描くプログラムを作成せよ。
// 形式：左端に値を表示し、適切に空白を空けて":"を書く（:で揃えるためにタブ\tを使うとよい）。その後ろに値の数だけ*を描くが、5個おきに空白を１つ入れる。具体例は下記の実行例を参照すること。
// * ヒント：入力値は配列に格納する。グラフを描くためには、値のループと、*を描くループの二重ループとなる。
package main

import "fmt"

func main() {
	counter := 5
	data := make([]int, counter)
	for i := 0; i < counter; i++ {
		fmt.Printf("input data[%v]: ", i)
		fmt.Scan(&(data[i]))
	}

	for i := 0; i < counter; i++ {
		fmt.Printf("%v\t:", data[i])
		for j := 1; j <= data[i]; j++ {
			fmt.Print("*")
			if (j % 5) == 0 {
				fmt.Print(" ")
			}
		}
		fmt.Println("")
	}
}
