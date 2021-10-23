package main

import (
	"fmt"
	s "strings"
)

// あとで何回も使うので、fmt.Println に短い名前のエイリアスを定義します。
var p = fmt.Println

func main() {
	// strings パッケージで使える関数の例は次の通りです。 これらは
	// パッケージで定義される関数であって、 文字列オブジェクトそのもの
	// のメソッドではないため、 関数への第 1 引数として対象の文字列を
	// 渡す必要があります。 さらに多くの関数を strings パッケージの
	// ドキュメントで確認できます。
	p("Contains:  ", s.Contains("test", "es"))
	p("Count:     ", s.Count("test", "t"))
	p("HasPrefix: ", s.HasPrefix("test", "te"))
	p("HasSuffix: ", s.HasSuffix("test", "st"))
	p("Index:     ", s.Index("test", "e"))
	p("Join:      ", s.Join([]string{"a", "b"}, "-"))
	p("Repeat:    ", s.Repeat("a", 5))
	p("Replace:   ", s.Replace("foo", "o", "0", -1))
	p("Replace:   ", s.Replace("foo", "o", "0", 1))
	p("Split:     ", s.Split("a-b-c-d-e", "-"))
	p("ToLower:   ", s.ToLower("TEST"))
	p("ToUpper:   ", s.ToUpper("test"))
	p()

	// strings パッケージには含まれませんが、 文字列のバイト長や、
	// 指定インデックスのバイトを取得する方法にここで触れておくことは
	// 意味があるでしょう。
	p("Len: ", len("hello"))
	p("Char:", "hello"[1])

	/*
	   len 関数やインデックス指定は、 バイトレベルで動作することに注意
	   してください。 Go は UTF-8 でエンコードされた文字列を使うので、
	   大抵そのままで役に立ちます。 しかし、もしマルチバイト文字列を扱う
	   可能性があるなら、 エンコードを考慮した操作をしたいと思うでしょう。
	   さらなる情報のために、"strings, bytes, runes and characters
	    in Go" を参照してください。
	*/
}
