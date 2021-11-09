package main

import "fmt"

// この intSeq 関数は、 intSeq の中で定義した
// 無名関数を返します。返された関数は、変数 i を
// 閉じ(closes over)てクロージャを作ります。
func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {
	// intSeq を呼び出した結果 (関数) を nextInt
	// に代入します。この関数は、自分自身の i の値を
	// キャプチャし、その値は nextInt を呼び出すた
	// びに更新されます。
	nextInt := intSeq()

	// nextInt を何回か呼び出して、クロージャの
	// 効果を見てください。
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	// 関数ごとに個別の状態をもっていることを確認する
	// ために、もう1つ新たに作成して試してみてください。
	newInts := intSeq()
	fmt.Println(newInts())
}
