package main

import "fmt"

func main() {
	// queue チャネルの2つの値を反復処理するとします。
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	// この range は、queue から受信した要素を反復処理します。
	// 上でチャネルを close したので、反復処理は2つの要素を
	// 受信したときに終了します。
	for elem := range queue {
		fmt.Println(elem)
	}
}
