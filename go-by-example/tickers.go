package main

import (
	"fmt"
	"time"
)

func main() {
	// ティッカーは、タイマーと同様の仕組み、すなわち値を送信するチャネルを
	// 使います。 ここでは、500ミリ秒ごとに受信される値を待つために、
	// select ビルトイン関数を使っています。
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	// ティッカーは、タイマーと同じように停止できます。 ティッカーが
	// 停止されると、そのチャネルからはもう 値を受信しなくなります。
	// この例では、1600ミリ秒後に停止します。
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}
