package main

// この構文は、encoding/base64 パッケージをデフォルトの
// base64 の代わりに b64 という名前でインポートします。
// サンプルコードのスペースを多少省略できます。
import (
	b64 "encoding/base64"
	"fmt"
)

func main() {
	// エンコード/デコードする string は次の通りです。
	data := "abc123!?$*&()'-=@~"

	// Go 標準の Base64 と URL 互換の Base64 の両方を
	// サポートしています。 以下は、標準のエンコーダーを使
	// ってエンコードする方法です。 エンコーダーは []byte
	// を要求するので、 string を型変換してやります。
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)

	// デコードはエラーを返す可能性があります。入力が適切か
	// あらかじめ分からない場合にはチェックできます。
	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	fmt.Println(string(sDec))
	fmt.Println()

	// 次のエンコード/デコードは、URL 互換の Base64
	// フォーマットを使います。
	uEnc := b64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println(uEnc)
	uDec, _ := b64.URLEncoding.DecodeString(uEnc)
	fmt.Println(string(uDec))
}
