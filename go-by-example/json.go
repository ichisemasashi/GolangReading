package main

import (
	"encoding/json"
	"fmt"
)

/*
   これら 2 つの構造体は、カスタムデータ型のエンコードと
   デコードをデモするために使います。
*/
type response1 struct {
	Page   int
	Fruits []string
}

/*
   エクスポートされたフィールドだけが JSON にエンコード/デコード
   されます。 そのため、フィールド名は大文字から始まる必要があります。
*/
type response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func main() {
	// まず最初に、基本的なデータ型から JSON 文字列への
	// エンコードを確認します。アトミックな値の例は次の通りです。
	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))

	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))

	fltB, _ := json.Marshal(2.34)
	fmt.Println(string(fltB))

	strB, _ := json.Marshal("gopher")
	fmt.Println(string(strB))

	// スライスとマップの場合は次のようになります。期待通り、
	// JSON の配列とオブジェクトにエンコードされます。
	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]int{"apple": 5, "lettuce": 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	// JSON パッケージは、カスタムデータ型も自動的に エンコードできます。
	// エンコード結果には、 エクスポートされたフィールドだけを含み、
	// デフォルトではそのフィールド名が JSON キーになります。
	res1D := &response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"},
	}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	// エンコードされる JSON キー名をカスタマイズするには、
	// 構造体のフィールド宣言にタグを指定します。 response2
	// の定義を確認してみてください。
	res2D := &response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"},
	}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))
}
