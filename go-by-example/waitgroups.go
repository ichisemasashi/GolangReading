package main

import (
	"fmt"
	"sync"
	"time"
)

/*
   これは、ゴルーチンで実行する関数です。 WaitGroup はポインタで
   関数に渡さなければならない点に注意してください。
*/
func worker(id int, wg *sync.WaitGroup) {
	// return するときに、完了したことを WaitGroup に通知します。
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)
	// 時間のかかるタスクをシミュレートするためにスリープします。
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	// この WaitGroup は、ここで起動したすべての
	// ゴルーチンが完了するのを待つために使われます。
	var wg sync.WaitGroup

	// いくつかのゴルーチンを起動し、そのたびに
	// WaitGroup のカウンターをインクリメントします。
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	// WaitGroup のカウンターが 0 に戻るまで（つまり、すべての
	// ワーカーが完了したと通知してくるまで）ブロックします。
	wg.Wait()
}
