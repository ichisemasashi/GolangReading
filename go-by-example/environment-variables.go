package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// キーと値のペアを設定するには、os.Setenv を使います。
	// キーに対する値を取得するには、os.Getenv を使います。
	// 環境にキーが存在しなければ、空文字列が返されます。
	os.Setenv("FOO", "1")
	fmt.Println("FOO:", os.Getenv("FOO"))
	fmt.Println("BAR:", os.Getenv("BAR"))

	// 環境に定義されたすべてのキーと値のペアを列挙するには、
	// os.Environ を使います。これは、KEY=value という形式
	// の文字列のスライスを返します。キーと値をそれぞれ取得する
	// ために、strings.SplitN を使えます。次の例は、すべての
	// キーを出力します。
	fmt.Println()
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0])
	}
}
