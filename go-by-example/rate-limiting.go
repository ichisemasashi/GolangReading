package main

import (
	"fmt"
	"time"
)

func main() {
	// まず最初に、基本的なレートリミットを見ていきます。リクエストに
	// 対する処理を制限したいとしましょう。これらのリクエストに同じ名前
	// の requests チャネルで対応します。
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	// limiter チャネルは、200ミリ秒ごとに値を受信します。これは、
	// レートリミットの仕組みの中でレギュレーターの役割を果たします。
	limiter := time.Tick(200 * time.Millisecond)

	// 各リクエストを処理する前に limiter チャネルからの受信を
	// ブロックさせることで、200ミリ秒に1リクエストしか処理できない
	// よう制限しています。
	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	// あるいは、全体としてはレートリミットを担保しつつ、一時的な
	// バーストリクエストは許容したいと思うかもしれません。
	// その場合、limiter チャネルをバッファリングすれば実現できます。
	// ここで、burstyLimiter チャネルは、3リクエストまでなら
	// バーストを許します。
	burstyLimiter := make(chan time.Time, 3)

	// 許容されているバースト量を表すため、チャネルに値を満たしておきます。
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}
	// 200 ミリ秒ごとに、burstyLimiter の上限である3つまで、
	// 新しい値を追加します。
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	// それでは、5リクエストきた場合をシミュレートします。
	// これらのうち最初の3リクエストは、burstyLimiter の
	// バースト機能の恩恵を受けるはずです。
	bustyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		bustyRequests <- i
	}
	close(bustyRequests)
	for req := range bustyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}
