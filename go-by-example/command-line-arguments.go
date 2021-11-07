package main

import (
	"fmt"
	"os"
)

func main() {
	// os.Args は、生のコマンドライン引数へのアクセスを提供します。
	// このスライスの 1 つ目の値は、プログラム自身のパスで、
	// os.Args[1:] がプログラムへの引数を保持することに注意してください。
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	// 個々の引数はインデックスを使って通常通り取得できます。
	arg := os.Args[3]

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)
}
