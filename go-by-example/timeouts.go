package main

import (
	"fmt"
	"time"
)

func main() {
	// 例として、2 秒後にチャネル c1 へ結果を返す外部呼び出しを実行していると仮定
	// しましょう。 このチャネルはバッファリングされるので、ゴルーチン内の送信は
	// ブロックしないことに注意してください。これは、チャネルが受信されない場合に
	// ゴルーチンのリークを防ぐ一般的な方法です。
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	// select を使ったタイムアウトの実装は次の通りです。 res := <-c1 が結果を
	// 待ち、<-time.After は 1 秒のタイムアウト後に送信されてくる値を待ちます。
	// select は最初に受信したものを処理するので、操作が 1 秒以上かかるとタイム
	// アウトのケースが選択されます。
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout 1")
	}

	// もしタイムアウトをさらに長い 3 秒にすると、 c2 からの受信が先に成功し、
	// 結果が表示されます。
	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}
}
