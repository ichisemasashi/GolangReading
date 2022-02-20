// No. 54 最大最小
// まずデータの個数を入力させ、次にデータの個数だけ整数値を入力させる。この入力データの中で最大値と最小値を求め表示するプログラムを作成せよ。データの個数は100個までとする。なお、データの個数とデータはファイルからリダイレクトで入力させればよいので、入力のためのメッセージは不要である（実行例を参照すること）。
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	var data []int

	data = readLines(data)

	max, min := max_min(data)
	fmt.Printf("最小値 = %v, 最大値 = %v\n", min, max)
}

func readLines(s []int) []int {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		x, _ := strconv.Atoi(scanner.Text())
		s = append(s, x)
	}
	return s
}

func max_min(sl []int) (int, int) {
	sort.Ints(sl)
	max := sl[len(sl)-1]
	min := sl[0]
	return max, min
}
