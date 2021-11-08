package main

import (
	"fmt"
	"os"
)

func main() {
	// os.Exit を使う場合は、defer は実行されません 。
	// そのため、この fmt.Println は決して呼ばれません。
	defer fmt.Println("!")

	// ステータス 3 で終了します。
	os.Exit(3)
}
