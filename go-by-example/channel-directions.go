package main

import "fmt"

/*
   この ping 関数は、チャネルを送信専用で受け取ります。このチャネルで
   受信しようとすると、コンパイルエラーになります。
*/
func ping(pings chan<- string, msg string) {
	pings <- msg
}

/*
   この pong 関数は、1つ目のチャネルを受信専用で (pings)、
   2つ目のチャネルを送信専用で (pongs) 受け取ります。
*/
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
