package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Go のシグナル通知は、チャネルに os.Signal 値を送信することで
	// 行います。 これらの通知を受信するためのチャネル(と、プログラムが
	// 終了できることを通知するためのチャネル) を作ります。
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// signal.Notify は、指定されたシグナル通知を受信するために、
	// 与えられたチャネルを登録します。
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// このゴルーチンは、シグナルを同期的に受信します。シグナルを受信し
	// たら、それを表示して、プログラムに終了できることを通知します。
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	// プログラムはシグナルを受信するまで (前述の done に値を送信
	// するゴルーチンで知らされる) 待機した後、終了します。
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
