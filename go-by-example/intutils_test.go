// デモのため、このコードは main パッケージにありますが、
// 任意のパッケージにできます。通常、テストコードはテスト
// されるコードと同じパッケージに配置します。
package main

import (
	"fmt"
	"testing"
)

// この整数の最小値を求めるシンプルな実装をテストしたい
// としましょう。 通常、テストしたいコードは intutils.go
// のような名前の ソースファイル中にあり、そのテストファイル
// は intutils_test.go のような名前になります。
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// テストは、Test から始まる名前の関数を書くことで作ります。
func TestIntMinBasic(t *testing.T) {
	ans := IntMin(2, -2)
	if ans != -2 {
		// t.Error* はテストの失敗を報告しますが、テストの
		// 実行は継続します。 t.Fatal* はテストの失敗を報告し、
		// テストの実行を即座に停止します。
		t.Errorf("IntMin(2, -2) = %d; want -2", ans)
	}
}

// テストを書くことは繰り返しになりがちなので、イディオムとして
// テーブル駆動スタイル (table-driven style) があります。
// この方法は、テストの入力と出力の期待値をテーブルに列挙し、
// 1レコードずつ順番にループしてロジックをテストします。
func TestIntMinTableDriven(t *testing.T) {
	var tests = []struct {
		a, b int
		want int
	}{
		{0, 1, 0},
		{1, 0, 0},
		{2, -2, -2},
		{0, -1, -1},
		{-1, 0, -1},
	}

	for _, tt := range tests {
		// t.Run はテーブルのエントリ1つずつに対して “サブテスト”
		// の実行を実現します。これらは、go test -v を実行したとき
		// に 個別に表示されます。
		testname := fmt.Sprintf("%d,%d", tt.a, tt.b)
		t.Run(testname, func(t *testing.T) {
			ans := IntMin(tt.a, tt.b)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
