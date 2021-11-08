package main

import "fmt"

func main() {
	// 空のマップを作成するには、make ビルトイン関数を
	// 使って make(map[キーの型]値の型) とします。
	m := make(map[string]int)

	// 典型的な name[key] = val という記法を使って、
	// キーと値のペアを設定します。
	m["k1"] = 7
	m["k2"] = 13
	// 例えば、 fmt.Println などでマップを表示すると、
	// すべてのキーと値のペアが出力されるでしょう。
	fmt.Println("map:", m)

	// name[key] でキーに対する値を取得できます。
	v1 := m["k1"]
	fmt.Println("v1: ", v1)

	// len ビルトイン関数は、マップに対しては
	// キーと値のペアの数を返します。
	fmt.Println("len:", len(m))

	// delete ビルトイン関数は、マップから
	// キーと値のペアを削除します。
	delete(m, "k2")
	fmt.Println("map:", m)

	// マップから値を取得するときの戻り値にはオプションで
	// 2 番目があり、指定したキーがマップに存在したかどうかを
	// 表します。これは、キーが存在しなかったのか、 0 や ""
	// などのゼロ値で存在したのかを区別するのに使えます。
	// この例では、値自体は必要なかったので、ブランク識別子
	// (blank identifier) と呼ばれる _ で無視しています。
	_, prs := m["k2"]
	fmt.Println("prs:", prs)

	// 新しいマップの宣言と初期化を次のように
	// 同じ行で書くこともできます。
	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("map:", n)
}
