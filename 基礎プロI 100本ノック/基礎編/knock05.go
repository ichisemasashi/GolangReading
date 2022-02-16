// No. 05 四則演算
// 整数値を2つ入力させ、それらの値の和、差、積、商と余りを求めるプログラムを作成せよ。なお、差と商は1つ目の値から2つ目の値を引いた、あるいは割った結果とする。余りは無い場合も0と表示するのでよい。
package main

import "fmt"

func main() {
	x := readNumber("input 1st number: ")
	y := readNumber("input 2nd number: ")
	printAns(x, y)
}

func readNumber(prom string) int {
	fmt.Print(prom)
	var n int
	fmt.Scan(&n)
	return n
}

func printAns(x, y int) {
	fmt.Println("和:", x+y)
	fmt.Println("差:", x-y)
	fmt.Println("積:", x*y)
	fmt.Printf("商: %v, 余り: %v\n", x/y, x%y)
}
