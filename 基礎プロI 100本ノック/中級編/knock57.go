// No. 57 テスト集計
// まず受験者数を入力させ、次に受験者数ごとに英語、数学、国語の点数をスペースで区切って入力させる（具体的な入力方法は下記のscanfの使い方の説明、および入力データの中身を見よ）。入力が終了したら、英語、数学、国語の平均点、および各受験生の合計点を計算して表示するプログラムを作成せよ。受験者数は100人までとする。なお、データの個数とデータはファイルからリダイレクトで入力させればよいので、入力のためのメッセージは不要である（実行例を参照すること）。
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type tensu struct {
	english int
	math    int
	lang    int
}

func main() {
	data, lines := readData()
	sums := make([]int, lines)
	avr := tensu{0, 0, 0}
	for i := 0; i < lines; i++ {
		avr.english += data[i].english
		avr.math += data[i].math
		avr.lang += data[i].lang
		sums[i] = data[i].english + data[i].math + data[i].lang
	}
	avr.english /= lines
	avr.math /= lines
	avr.lang /= lines

	fmt.Printf("平均点 英語:%v, 数学:%v, 国語:%v\n", avr.english, avr.math, avr.lang)
	fmt.Println("個人合計点")
	for i := 0; i < lines; i++ {
		fmt.Printf("%v: %v\n", i, sums[i])
	}
}

func readData() ([]tensu, int) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	lines, _ := strconv.Atoi(scanner.Text())

	data := make([]tensu, lines)
	for i := 0; i < lines; i++ {
		scanner.Scan()
		data[i].english, _ = strconv.Atoi(scanner.Text())
		scanner.Scan()
		data[i].math, _ = strconv.Atoi(scanner.Text())
		scanner.Scan()
		data[i].lang, _ = strconv.Atoi(scanner.Text())
	}
	return data, lines
}
