// No. 59 行列の和
// 3x3行列の和を求めて表示するプログラムを作成せよ。行列の値は2次元配列で表現し、繰り返しを使って計算すること。
// 3x3行列とは縦3つ、横3つの9つの要素(値)をひとまとめにして扱うものである。2つの3x3行列の和は次式のように、それぞれ同じ位置にある値を足したものとして計算できる。
// 例えばa12という要素は、1行目2列目の要素という意味である。それぞれ同じ位置にある要素を足せばよい。
// なお、入力値は1行ずつ3つの値をスペースで区切って入力するようにするとよい。このためには、scanf("%d %d %d", &a[0][0], &a[0][1], &a[0][2]);のように書く(No. 57参照)。
package main

import "fmt"

func main() {
	a1 := [3][3]int{}
	a2 := [3][3]int{}
	fmt.Println("1つめの行列")
	for i := 0; i < 3; i++ {
		fmt.Scanf("%d %d %d", &a1[i][0], &a1[i][1], &a1[i][2])
	}
	fmt.Println("2つめの行列")
	for i := 0; i < 3; i++ {
		fmt.Scanf("%d %d %d", &a2[i][0], &a2[i][1], &a2[i][2])
	}

	fmt.Println("和")
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Print(a1[i][j]+a2[i][j], "\t")
		}
		fmt.Println("")
	}
}
