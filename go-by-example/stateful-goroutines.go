package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

/*
   この例では、状態は1つのゴルーチンに所有させます。これは、データが
   並行アクセスで決して壊れないことを保証します。ほかのゴルーチンが
   状態を読み書きするためには、データを所有するゴルーチンにメッセージ
   を送り、対応するデータを受け取る必要があります。 readOp 構造体と
    writeOp 構造体は、それらの要求やデータをもつゴルーチンが応答する
	方法をカプセル化します。
*/
type readOp struct {
	key  int
	resp chan int
}

type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func main() {
	// 前回と同様に、操作回数をカウントします。
	var readOps uint64
	var writeOps uint64
	// reads および writes チャネルは、ほかのゴルーチンから
	// 読み書きのリクエストを発行するために使われます。
	reads := make(chan readOp)
	writes := make(chan writeOp)

	// これは前回の例と同様にマップ state を持つゴルーチンですが、
	// ゴルーチンにプライベートな変数となっている点が異なります。
	// このゴルーチンは、reads チャネルと writes チャネルを繰り返し
	// select し、リクエストが届いたら応答します。レスポンスは、
	// まず要求された操作を実行し、その後成功を示す値 (reads の場合は
	// 要求された値) を応答チャネル resp に送信します。
	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	// ここで、状態をもつゴルーチンから reads チャネル経由で
	// 値を読み込むゴルーチンを100個開始します。各読み込みは、
	// readOp を構築し, reads チャネルへ送信し、結果を resp
	// チャネルから受信する必要があります。
	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := readOp{
					key:  rand.Intn(5),
					resp: make(chan int),
				}
				reads <- read
				<-read.resp
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// 同様に、書き込み用のゴルーチンを10個開始します。
	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool),
				}
				writes <- write
				<-write.resp
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// ゴルーチンを1秒間動作させます。
	time.Sleep(time.Second)

	// 最後に、操作回数を取得してレポートします。
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps", writeOpsFinal)
}
