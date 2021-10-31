package main

import (
	"fmt"
	"time"
)

func main() {
	p := fmt.Println

	// これは RFC3339 に対応するレイアウト定数を使って、
	// 日時をフォーマットする基本的な例です。
	t := time.Now()
	p(t.Format(time.RFC3339))

	// 日時をパースする場合も Format と同じレイアウト値を使います。
	t1, e := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	p(t1)

	// Format と Parse は、サンプルベースのレイアウトを使います。
	// 通常は、time パッケージに定義された定数を使いますが、
	// カスタムレイアウトを指定することもできます。 レイアウトは、
	// フォーマットやパースに使うパターンを示すために、 参照日時である
	// Mon Jan 2 15:04:05 MST 2006 を使わなければなりません。
	// サンプル日時は、年が 2006、時が 15、曜日が月曜、 とまさに
	// この通りの値を指定しなければなりません。
	p(t.Format("3:04PM"))
	p(t.Format("Mon Jan _2 15:04:05 2006"))
	p(t.Format("2006-01-02T15:04:05.999999-07:00"))
	form := "3 04 PM"
	t2, e := time.Parse(form, "8 41 PM")
	p(t2)

	// 単なる数値表現だけであれば、標準的な文字列フォーマットに
	// 日時を構成する個々の値を抽出して指定することもできます。
	fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	// Parse は、不正な入力に対しては、問題を説明するエラーを返します。
	ansic := "Mon Jan _2 15:04:05 2006"
	_, e = time.Parse(ansic, "8:41PM")
	p(e)
}
