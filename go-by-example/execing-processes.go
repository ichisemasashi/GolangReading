package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// 例として ls を実行します。 Go は、実行したいバイナリの
	// 絶対パスを要求するので、 exec.LookPath を使って探します
	// (おそらく /bin/ls)。
	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}

	// Exec は、引数を (1 つの文字列ではなく) スライスで要求
	// します。いくつか一般的な引数を ls に渡してみることにし
	// ましょう。最初の引数は、プログラムの名前でなければならない
	// 点に注意してください。
	args := []string{"ls", "-a", "-l", "-h"}

	// また、Exec には使用する 環境変数 も必要です。ここでは、
	// 現在の環境変数をそのまま渡すことにします。
	env := os.Environ()

	// これが、実際の syscall.Exec 呼び出しです。この呼び出し
	// が成功すると、私たちのプロセスの実行は終了し、
	//  /bin/ls -a -l -h プロセスで置き換わります。
	// もしエラーが発生すると、戻り値を受け取ることになります。
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}
