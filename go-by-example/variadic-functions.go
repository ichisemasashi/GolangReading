package main

import "fmt"

// これは任意個の int を引数として受け取る関数です。
func sum(nums ...int) {
	fmt.Print(nums, " ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

func main() {
	// 可変長引数関数は、通常通り個々の
	// 引数を渡して呼び出せます。
	sum(1, 2)
	sum(1, 2, 3, 4)

	// 複数の引数をすでにスライスでもっている場合は、
	// func(slice...) のような形で可変長引数関数に渡せます。
	nums := []int{1, 2, 3, 4}
	sum(nums...)
}
