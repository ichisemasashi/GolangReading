package main

import "fmt"

func main() {
	// この例では、string を 2 つまでバッファリングするチャネルを make しています。
	messages := make(chan string, 2)

	// このチャネルはバッファリングされるので、対応する受信側が
	// いなくても値をチャネルに送信できます。
	messages <- "buffered"
	messages <- "channel"

	// そして、あとで通常通り 2 つの値を受信できます。
	fmt.Println(<-messages)
	fmt.Println(<-messages)
}
