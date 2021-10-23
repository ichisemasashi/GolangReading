package main

import (
	"fmt"
	"sort"
)

/*
   Go でカスタム関数を使ってソートするためには、 対応する型が必要です。
   ここでは、byLength 型を作りました。これは、 組み込みの []string 型の
   ただのエイリアスです。
*/
type byLength []string

/*
   sort パッケージの Sort 関数を使えるように、 sort.Interface
   すなわち Len, Less, Swap 関数を実装します。 Len と Swap は
   どの型でもだいたい同じになり、 Less が実際のカスタムソートのロジックを
   もちます。 この例では、文字列の長さの昇順でソートしたいので、
   len(s[i]) と len(s[j]) を使っています。
*/
func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func main() {
	// 元の fruits スライスを byLength に型変換し、
	// sort.Sort 関数を使うことでカスタムソートを実現できます。
	fruits := []string{"peach", "banana", "kiwi"}
	sort.Sort(byLength(fruits))
	fmt.Println(fruits)
}
