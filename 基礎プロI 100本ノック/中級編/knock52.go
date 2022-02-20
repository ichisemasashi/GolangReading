// No. 52 閏年
// 西暦を入力したらその年が閏（うるう）年かどうかを判定するプログラムを作成せよ。なお、4で割り切れる年のうち、100で割り切れないか、400で割り切れる年は閏年である。
package main

import "fmt"

func main() {
	var y int
	fmt.Print("input year: ")
	fmt.Scan(&y)

	if (y % 4) == 0 {
		if ((y % 400) == 0) || ((y % 100) != 0) {
			fmt.Printf("%vは閏年である\n", y)
		} else {
			fmt.Printf("%vは閏年でない\n", y)
		}
	} else {
		fmt.Printf("%vは閏年でない\n", y)
	}
}
