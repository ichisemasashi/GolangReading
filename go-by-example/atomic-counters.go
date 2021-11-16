package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// 正数カウンターのために符号なし整数を使います。
	var ops uint64
	// WaitGroup は、すべてのゴルーチンがタスク
	// を完了するのを待つときに使えます。
	var wg sync.WaitGroup

	// カウンターをちょうど1000回インクリメント
	// するゴルーチンを50個開始します。
	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			for c := 0; c < 1000; c++ {
				// カウンターをアトミックにインクリメントするため、
				// & 構文で ops カウンターのメモリアドレスを
				// AddUint64 に与えます。
				atomic.AddUint64(&ops, 1)
			}
			wg.Done()
		}()
	}

	// すべてのゴルーチンが完了するまで待ちます。
	wg.Wait()

	// 書き込み中のゴルーチンがないことを知っているので、
	// ops に安全にアクセスできます。atomic.LoadUint64 のような
	// 関数を使えば、更新され続けているカウンターを安全に
	// 読み込むことも可能です。
	fmt.Println("ops:", ops)
}
