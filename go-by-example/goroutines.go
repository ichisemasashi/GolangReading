package main

import (
	"fmt"
	"time"
)

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

/* このプログラムを実行すると、最初に同期呼び出しの出力、 その次に
   2 つのゴルーチンの混ざった出力を確認できます。 これは、ゴルーチンが
   Go ランタイムによって 並行実行されていることを示しています。 */
func main() {
	// 関数呼び出し f(s) があるとしましょう。
	// 通常の方法で呼び出すと、同期的に実行されます。
	f("direct")

	// この関数をゴルーチンとして呼び出すには、 go f(s) とします。
	// この新しいゴルーチンは、 実行元とは並行して実行されます。
	go f("goroutine")

	// 無名関数に対してゴルーチンを開始することもできます。
	go func(msg string) {
		fmt.Println(msg)
	}("going")

	// 上の 2 つの関数呼び出しは別々のゴルーチンで非同期に実行
	// されるので、プログラムの実行はすぐにここへきます。 それらが
	// 完了するまで待ちます（より堅牢な方法としては、 WaitGroup を
	// 使ってください）。
	time.Sleep(time.Second)
	fmt.Println("done")
}
