package main

// Go は、基本的なコマンドラインフラグのパースをサポートする
// flag パッケージを提供しています。 サンプルのコマンドライン
// プログラムを実装するために、 このパッケージを使いましょう。
import (
	"flag"
	"fmt"
)

func main() {
	// 基本的なフラグ宣言は、文字列、整数、真偽値に対して
	// 利用できます。ここでは、文字列フラグ word をデフォルト値
	// "foo" と簡単な説明とともに宣言します。 この flag.String
	// 関数は、(文字列値ではなく) 文字列ポインタを返します。
	// このポインタの使い方は後ほど確認します。
	wordPtr := flag.String("word", "foo", "a string")

	// 次の例は、word フラグと同じ方法を使って、 numb フラグと
	// fork フラグを宣言しています。
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	// プログラムの別の場所で宣言された変数を使うようなオプションを
	// 宣言することも可能です。フラグを宣言する関数にポインタの形で
	// 渡す必要がある点に 注意してください。
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	// すべてのフラグを宣言したら、コマンドラインをパースするために
	// flag.Parse() を呼び出します。
	flag.Parse()

	// ここでは、パースしたオプションと残りの任意の引数を単に出力します。
	// 実際のオプション値を取得するために、*wordPtr のようにポインタを
	// デリファレンスする必要がある点に気をつけてください。
	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
}
