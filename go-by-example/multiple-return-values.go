package main

import "fmt"

// この関数シグネチャの (int, int) は、この関数が
// 2 つの int を返すことを示しています。
func vals() (int, int) {
	return 3, 7
}

func main() {
	// ここでは、 多重代入 (multiple assignment)
	// で関数呼び出しからの 2 つの戻り値を使っています。
	a, b := vals()
	fmt.Println(a)
	fmt.Println(b)

	// もし戻り値の一部しか欲しくないのであれば、
	// ブランク識別子 _ を使ってください。
	_, c := vals()
	fmt.Println(c)
}
