package main

import (
	"fmt"
	"time"
)

/*
   ゴルーチンで実行する関数は次の通りです。この関数が完了したことを
   別のゴルーチンに通知するため、 done チャネルが使われます。
*/
func woker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	// 完了したことを通知するために値を送信します。
	done <- true
}

/*
   このプログラムから <-done の行を削除すると、 worker が
   開始さえする前にプログラムが終了してしまいます。
*/
func main() {
	// 通知用のチャネルを渡して、worker ゴルーチンを開始します。
	done := make(chan bool, 1)
	go woker(done)

	// チャネルへの完了通知を受信するまでブロックします。
	<-done
}
