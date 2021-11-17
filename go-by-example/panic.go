package main

import "os"

func main() {
	// 予期せぬエラーを確認するために、panic を使います。
	// これは、panic するために作られた唯一のプログラムです。
	panic("a problem")

	// panic の一般的な使い方は、ある関数が扱い方を知らない
	// (または、扱いたくない) エラー値を返したときに、異常終了させる
	// ことです。新規ファイル作成時に予期せぬエラーが発生したら
	// panic する例は、次の通りです。
	_, err := os.Create("/tmp/file")
	if err != nil {
		panic(err)
	}
}
