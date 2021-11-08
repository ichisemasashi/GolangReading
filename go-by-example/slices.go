package main

import "fmt"

func main() {
	// 配列とは異なり、スライスは含まれる要素だけで 型が定義されます
	// (要素の数は関係ありません)。 長さが 0 でない空のスライスを作る
	// には、 make ビルトイン関数を使います。ここでは、長さ3の string
	// のスライスを作っています (ゼロ値で初期化されます)。
	s := make([]string, 3)
	fmt.Println("emp:", s)

	// 配列と全く同じように値の設定と取得が可能です。
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("set:", s)
	fmt.Println("get:", s[2])

	// len は期待通りスライスの長さを返します。
	fmt.Println("len:", len(s))

	// これらの基本的な操作に加えて、スライスは配列よりも便利な
	// 操作をサポートします。 その 1 つが、1 つ以上の新しい値を
	// 含んだ スライスを返す append ビルトイン関数です。 ただし、
	// append から返される値として 新しいスライス値を受け取るかも
	// しれないことを 許容する必要がある点に注意してください。
	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("apd:", s)

	// スライスは copy することもできます。ここでは、
	// s と同じ長さの空のスライス c を作成し、
	// s から c へコピーしています。
	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("cpy:", c)

	// スライスは、 slice[下限:上限] という記法の “スライス” 演算
	// をサポートします。例えば、これは s[2] と s[3], s[4] の要素
	// をもつスライスを 取得します。
	l := s[2:5]
	fmt.Println("sl1:", l)

	// これは、 s[5] まで (上限は除いて) スライスします。
	l = s[:5]
	fmt.Println("sl2:", l)

	// そしてこれは、s[2]から（下限は含んで）スライスします。
	l = s[2:]
	fmt.Println("sl3:", l)

	// スライスの変数を 1 行で宣言かつ初期化することもできます。
	t := []string{"g", "h", "i"}
	fmt.Println("dcl:", t)

	// スライスを多次元のデータ構造にすることもできます。
	// 多次元配列とは異なり、内側のスライスの長さは変わりえます。
	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}
