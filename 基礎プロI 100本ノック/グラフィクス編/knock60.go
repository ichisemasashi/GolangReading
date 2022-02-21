// No. 60 円
// 中心座標を入力させ、横600縦400のウィンドウを開き、入力した座標に半径50の円を描くプログラムを作成せよ。
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	win_Width = 600
	win_Hight = 400
	R         = 50
)

var circle_x, circle_y int

func CircleImage() image.Image {
	// 指定された幅と高さの新しい image.RGBA を作成し、その画像にレンダリングするためのコンテキストを準備します。
	dc := gg.NewContext(win_Width, win_Hight)

	// 現在の色を設定します。r、g、bの値は0から1の間でなければなりません。
	dc.SetRGB(0, 0, 0)
	dc.DrawCircle(float64(circle_x), float64(circle_y), R)

	// 現在のパスを、現在の色、線幅、ラインキャップ、ラインジョイン、ダッシュの設定でストロークします。この操作の後、パスはクリアされます。
	dc.Stroke()
	return dc.Image()
}

// Game は ebiten.Game インターフェースを実装しています。
type Game struct{}

// Updateはゲーム状態を進行させる。
// Updateはtick毎（デフォルトでは1/60 [s]）に呼ばれる。
func (g *Game) Update() error {
	return nil
}

// Drawはゲーム画面を描画します。
// Drawはフレーム毎（60Hz表示の場合は通常1/60[s]）に呼び出されます。
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White) // 背景色：白
	m := CircleImage()
	em := ebiten.NewImageFromImage(m)
	screen.DrawImage(em, &ebiten.DrawImageOptions{})
}

// Layoutは外部サイズ（ウィンドウサイズなど）を受け取り、（論理的な）画面サイズを返します。
// 外側のサイズで画面サイズを調整する必要がない場合は、固定サイズを返せばよい。
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return win_Width, win_Hight
}

func main() {
	fmt.Print("円の中心座標を入力: ")
	fmt.Scanf("%d %d", &circle_x, &circle_y)

	// ウィンドウサイズを任意に指定します。
	ebiten.SetWindowSize(win_Width, win_Hight)
	ebiten.SetWindowTitle("Hello, World!")
	// ebiten.RunGameを呼び出して、ゲームループを開始します。
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
