// No. 55 夢想花again
// 「とんで」を9回「まわって」を3回繰り返した後「まわる」と表示して改行する、を3回繰り返すプログラムを作成せよ。「とんで」「まわって」と3行文の繰り返しは必ず繰り返し構文を使うこと。
package main

import "fmt"

func main() {
	for i := 0; i < 3; i++ {
		repeatString("とんで", 9)
		repeatString("まわって", 3)
		repeatString("まわる", 1)
		fmt.Println("")
	}
}

func repeatString(s string, n int) {
	for i := 0; i < n; i++ {
		fmt.Print(s)
	}
}
