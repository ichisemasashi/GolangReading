// No. 45 タクシー料金
// 初乗り料金が1700mまで610円、それ以降は313mごとに80円のタクシーがある。このタクシーに乗った距離をm単位で入力し、料金を計算するプログラムを作成せよ。
// * 313mごとの区間は1mでも進んでしまったら80円かかることに注意。
package main

import "fmt"

func main() {
	var n int
	fmt.Print("距離: ")
	fmt.Scan(&n)

	pay := 0
	pay_base := 610
	pay_add := 80
	hatunori := 1700
	if n <= hatunori {
		pay = pay_base
	} else if ((n - hatunori) % 313) == 0 {
		pay = pay_base + int((n-hatunori)/313)*pay_add
	} else {
		pay = pay_base + int((n-hatunori)/313)*pay_add + pay_add
	}
	fmt.Println("金額", pay)
}
