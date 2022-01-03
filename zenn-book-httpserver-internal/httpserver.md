# Deep Dive into The Go\'s Web Server


Goのnet/httpパッケージはとてもよくできており、Webサーバーを動かすのに必要になる「httpコネクションを確立してリクエストを読んでルーティングして......」という手続き的な処理を気にせずとも誰でも簡単にWebサーバーを立てられるようになっています。
ですが、そのnet/httpが代わりにやってくれている「裏側の処理」の部分が気になる、何やっているんだろう？と不思議に思っている方はいませんか？
この本では、実際に筆者がnet/httpパッケージのソースコードを読み込んだうえで、「GoのWebサーバーがどのような仕組みで起動・動いているのか」というところについて、図を使いながら解説しています。



# はじめに

# この本について

Goの`net/http`パッケージはとても使い勝手がよく、コードを数行書いただけですぐにWebサーバーの\"Hello
World\"ができるようになります。


```
// Hello Worldの一部
h1 := func(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "Hello from a HandleFunc #1!\n")
}
http.HandleFunc("/", h1)
log.Fatal(http.ListenAndServe(":8080", nil))
```


自分が使いたいハンドラを書き、それをどこのURL・ポート番号で起動させるかということを**宣言的に**指定するだけで、簡単にその要件を満たしたWebサーバーが起動します。\
ですが、この裏に何が行われているのか考えてみたことはありますか？

筆者はこれをやったときに、

-   `net.Listen()`とか`ln.Accept()`とかやらなくていいの？
-   リクエストをネットワークから「Readする/読み込む」ような処理がどこにも書かれてないけど、どこでやってるの？
-   `http.HandleFunc`関数に登録しているハンドラって、どこに行くの？どこに保存されるの？
-   ハンドラに出てくる第一引数`http.ResponseWriter`って何者？なんでこれに書き込むだけで勝手にレスポンスになるの？
-   そもそも`http.ListenAndServe`の第二引数`nil`っていいの？

などというポイントが気になってしかたがありませんでした。~~我ながらめんどくさい性格してると思います~~

実際にGoのWebサーバーが行っている「コネクションを確立してリクエストを読み込んで......」という**手続き的な**処理については、`net/http`の中にまるっと隠蔽されユーザーが気にしなくていいようになっているので、それを知りたいなら実際にコードリーディングをするのが一番の近道でした。

この本では、当時私が抱いた疑問の答えを見つけるために`net/http`パッケージのソースコードを読み漁って得た内容を、わかりやすくまとめ図式化したものを共有するために執筆したものとなります。\
これを通して、「GoのWebサーバーの中身、ちょっとわかる」となってもらえたら嬉しいです。

## 本の構成

### 2章　ウェブサーバーのHello, World

この章では、「サーバーの中身の前に、どうやってGoでサーバーを立てるの？」というところを軽く紹介します。

### 3章　httpサーバー起動の裏側

ここでは、`http.ListenAndServe`関数が実行されてから、その第二引数に渡されたルーティングハンドラが起動するところまでのソースコードを追っていき、流れを図示しています。

### 4章 デフォルトでのルーティング処理の詳細

この章では、`ListenAndServe`関数の第二引数が`nil`であった場合、デフォルトで採用されるルータ`DefaultServeMux`の

-   ハンドラ登録
-   ハンドリング処理

の仕組みについて、`net/http`のソースコードを追って図示しています。

### 5章　ハンドラによるレスポンス返却の詳細

ここでは、ユーザーが用意したハンドラ関数`func(w http.ResponseWriter, _ *http.Request)`に関して、

-   第一引数`http.ResponseWriter`は何者なのか
-   `http.ResponseWriter`の書き込みメソッドを実行するだけで、ネットワークに載せるレスポンスが作れてしまうのはどういう仕組みなのか

という2点についてまとめました。

### 6章 GoでのWebサーバー起動の全体図

最後に、3\~5章に渡って説明してきた内容の要点だけを抽出し作成した「Goでウェブサーバーを起動したときに裏で起きている概要図」をここで紹介します。

## 使用する環境・バージョン

-   OS: macOS Catalina 10.15.7
-   go version go1.17 darwin/amd64

## 読者に要求する前提知識

-   Goの基本的な文法の読み書きができること




# ウェブサーバーのHello, World

# この章について

この章では「そもそもGoでWebサーバーを立てるにはどうしたらいい？」というHello
World部分を軽く紹介します。

# コード全容&デモ

`main.go`というファイル内に、以下のようなコードを用意します。


```
package main

import (
    "io"
    "log"
    "net/http"
)

func main() {
    h1 := func(w http.ResponseWriter, _ *http.Request) {
        io.WriteString(w, "Hello from a HandleFunc #1!\n")
    }
    h2 := func(w http.ResponseWriter, _ *http.Request) {
        io.WriteString(w, "Hello from a HandleFunc #2!\n")
    }

    http.HandleFunc("/", h1)
    http.HandleFunc("/endpoint", h2)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```


コード出典:pkg.go.dev -
net/http#example-HandleFunc

コードの内容については後ほど説明しますので、ひとまずこれを動かしてみましょう。\
ターミナルを開いて、`go run`コマンドを実行しましょう。


``` language-bash
$ go run main.go
```


すると、httpサーバーが`localhost:8080`で立ち上がります。

別のターミナルを開いて、`curl`コマンドを用いてリクエストを送信してみましょう。


``` language-bash
$ curl http://localhost:8080
Hello from a HandleFunc #1!

$ curl http://localhost:8080/endpoint
Hello from a HandleFunc #2!
```


このように、きちんとレスポンスを得ることができました。

# コード解説

それでは、先ほどの`main.go`の中では何をやっているのかについて見ていきましょう。

## 1. ハンドラを作成する


```
h1 := func(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "Hello from a HandleFunc #1!\n")
}
h2 := func(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "Hello from a HandleFunc #2!\n")
}
```


受け取ったリクエストに対応するレスポンスを返すためのハンドラ関数を作ります。\
Goではhttpハンドラは`func(w http.ResponseWriter, _ *http.Request)`の形で定義する必要があるので、そのシグネチャ通りの関数をいくつか作成します。

## 2. ハンドラとURLパスを紐付ける


```
http.HandleFunc("/", h1)
http.HandleFunc("/endpoint", h2)
```


`http.HandleFunc`関数に、「先ほど作ったハンドラは、どのURLパスにリクエストが来たときに使うのか」という対応関係を登録していきます。

## 3. サーバーを起動する


```
log.Fatal(http.ListenAndServe(":8080", nil))
```


ハンドラとパスの紐付けが終了したところで`http.ListenAndServe`関数を呼ぶと、今まで私たちが設定してきた通りのサーバーが起動されます。

# 次章予告

次の章からは、`http.ListenAndServe`関数を呼んだ後に、どのような処理を経てサーバーがリクエストを受けられるようになるのか、という内部実装部分を掘り下げていきます。




# httpサーバー起動の裏側

# この章について

Goでは`http.ListenAndServe`関数を呼ぶことで、httpサーバーを起動させることができます。


```
http.ListenAndServe(":8080", nil)
```


この章では、`http.ListenAndServe`が呼ばれた裏側で、どのような処理が行われているのかについて解説します。

# コードリーディング

Goの利点として「GoはGo自身で書かれているため、コードリーディングのハードルが低い」というのがあります。\
そのため、`net/http`パッケージに存在する`http.ListenAndServe`関数の実装コードももちろんGoで行われています。

ここからは、`http.ListenAndServe`関数の挙動を理解するために、`net/http`パッケージのコードを実際に読んでいきたいと思います。

## 1. `http.ListenAndServe`関数

`http.ListenAndServe`関数自体はとても単純な実装です。


```
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, handler Handler) error {
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}
```


出典:net/http/server.go

`http.Server`型を作成し、それの`ListenAndServe`メソッドを呼んでいることがわかります。\
![](https://storage.googleapis.com/zenn-user-upload/80b3cb5507083e854a422327.png)

このとき作られる`http.Server`型は、`http.ListenAndServe`関数の引数として渡された「サーバーアドレス」と「ルーティングハンドラ」を内部に持つことになります。


```
type Server struct {
    Addr string
    Handler Handler // handler to invoke, http.DefaultServeMux if nil
    // (以下略)
}
```


出典:pkg.go.dev - net/http#Server

サーバーアドレスとルーティングハンドラは、それぞれ`http.ListenAndServe`関数の第一引数、第二引数で指定されたものが使用されます。\
もし`http.ListenAndServe`の第二引数が`nil`だった場合は、`net/http`パッケージ内でデフォルトで用意されている`DefaultServeMux`が使用されます。

> `func ListenAndServe(addr string, handler Handler) error`\
> The `handler` is typically `nil`, in which case the `DefaultServeMux`
> is used.
>
> 出典:pkg.go.dev -
> net/http#ListenAndServe


```
// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux
```


出典:pkg.go.dev -
net/http#pkg-variables

## 2. `http.Server`型の`ListenAndServe`メソッド

`http.Server`型の`ListenAndServe`メソッドの中身は以下のようになっています。


```
func (srv *Server) ListenAndServe() error {
    if srv.shuttingDown() {
        return ErrServerClosed
    }
    addr := srv.Addr
    if addr == "" {
        addr = ":http"
    }
    ln, err := net.Listen("tcp", addr) // 1. net.Listenerを得る
    if err != nil {
        return err
    }
    return srv.Serve(ln) // 2. Serveメソッドを呼ぶ
}
```


出典:net/http/server.go

ここでやっていることは大きく2つです。

1.  `net.Listen`関数を使って、`net.Listener`インターフェース`ln`を得る
2.  `ln`を引数に使って、`http.Server`型の`Serve`メソッドを呼ぶ

![](https://storage.googleapis.com/zenn-user-upload/1c6e8adcfd5ee903c8e2116e.png)

## 3. `http.Server`型の`Serve`メソッド

次に、`http.Server`型の`ListenAndServe`メソッド中で呼ばれた`Serve`メソッドを見てみましょう。


```
func (srv *Server) Serve(l net.Listener) error {
    // (一部抜粋)
    // 1. contextを作る
    baseCtx := context.Background()
    ctx := context.WithValue(baseCtx, ServerContextKey, srv)

    for {
        rw, err := l.Accept() // 2. ln.Acceptをしてnet.Connを得る

        connCtx := ctx
        c := srv.newConn(rw) // 3. http.conn型を作る
        go c.serve(connCtx) // 4. http.conn.serveの実行
    }
}
```


出典:net/http/server.go

内部でやっているのは、以下の4つです。

1.  contextを作る
2.  `net.Listener`のメソッド`Accept()`を呼んで、`net.Conn`インターフェース`rw`を得る
3.  `net.Conn`から`http.conn`型を作る
4.  新しいゴールーチン上で、`http.conn`型の`serve`メソッドを実行する

![](https://storage.googleapis.com/zenn-user-upload/278785ca65cc5dcd5bf86e90.png)

この処理の中にはいくつか重要なポイントがありますので、ここからはそれを解説していきます。

### `net.Conn`インターフェースの入手

この時点で、http通信をするためのコネクションインターフェース`net.Conn`の入手が完了します。

`net.Conn`の入手のために必要なステップは2つです。

1.  (`(srv *Server) ListenAndServe`メソッド内)
    `net.Listen`関数から`net.Listener`インターフェース`ln`を得る
2.  (`(srv *Server) Serve`メソッド内) `ln.Accept()`メソッドを実行する


```
func (srv *Server) ListenAndServe() error {
    // (一部抜粋)
    ln, err := net.Listen("tcp", addr) // 1. net.Listenerを得る
    return srv.Serve(ln)
}

func (srv *Server) Serve(l net.Listener) error {
    // (一部抜粋)
    for {
        rw, err := l.Accept() // 2. ln.Acceptをしてnet.Connを得る
    }
}
```


`net.Conn`インターフェースには`Read`,`Write`メソッドが存在し、それらを実行することでネットワークからのリクエスト読み込み・レスポンス書き込みを行えるようになります。


`net.Conn`を利用したネットワークI/Oの詳細については、拙著Goから学ぶI/O
第3章をご覧ください。


### `for`無限ループによる処理の永続化

`ln.Accept()`メソッドによって得られた`net.Conn`は、一回の「リクエストーレスポンス」にしか使えません。\
つまりこれは、「一つの`net.Conn`を使い回す形で、サーバーにくる複数のリクエストを捌くことはできない」ということです。

そのため、`for`無限ループを利用して「一つのリクエストごとに一つの`net.Conn`を作成するのを繰り返す」ことでサーバーを継続的に稼働させているのです。


```
func (srv *Server) Serve(l net.Listener) error {
    // (一部抜粋)
    for {
        rw, err := l.Accept()
        go c.serve(connCtx)
    }
}
```


### 新規ゴールーチン上での`http.conn.serve`メソッド稼働

実際にリクエストをハンドルして、レスポンスを返す作業である`http.conn.serve`メソッドは、`http.ListenAndServe`関数が動いているメインゴールーチン上ではなく、`go`文によって作成される新規ゴールーチン上にて実行されています。


```
// (再掲)
func (srv *Server) Serve(l net.Listener) error {
    for {
        go c.serve(connCtx) // 4. http.conn.serveの実行
    }
}
```


わざわざ新規ゴールーチンを立てるのは、リクエストの処理を並行に実施できるようにするためです。

メインゴールーチン上でリクエストを逐次的に処理してしまうと、一つ時間がかかるリクエストが来た場合に、その間にきた別のリクエストはその時間がかかっているリクエスト処理が終わるまで待たされることになってしまいます。\
1リクエストごとに新規ゴールーチンを立てた場合は、複数リクエストを並行に処理できるようになるためレスポンスタイムが向上します。

![](https://storage.googleapis.com/zenn-user-upload/719f24627dd37426cbe28e76.png)

## 4. `http.conn`型の`serve`メソッド

本題の`http.ListenAndServe`関数の掘り下げに戻りましょう。\
`http.Server`型`Serve`メソッド内で、`http.conn.serve`メソッドが呼ばれたところまで見てきました。

ここからは`http.conn.serve`メソッドを見ていきます。


```
func (c *conn) serve(ctx context.Context) {
    // 一部抜粋
    for {
        w, err := c.readRequest(ctx)
        serverHandler{c.server}.ServeHTTP(w, w.req)
    }
}
```


出典:net/http/server.go

`http.conn.serve`内部で行っているのは、大きく分けて以下の2つです。

1.  `http.conn`型の`readRequest`メソッドから、`http.response`型を得る
2.  `http.serverHandler`型の`ServerHTTP`メソッドを呼ぶ

![](https://storage.googleapis.com/zenn-user-upload/95510eac8b4dacc3e2f7b19a.png)

これも一つずつ詳しく説明していきます。

### 4-1. `http.conn.readRequest`メソッドによる`http.response`型の入手

まず、`readRequest`型のレシーバである`http.conn`型は、内部に先ほど入手した`net.Conn`を含んでいます。


```
// A conn represents the server side of an HTTP connection.
type conn struct {
    server *Server
    rwc net.Conn
    // (以下略)
}
```


出典:net/http/server.go

この`net.Conn`の`Read`メソッドを駆使してリクエスト内容を読み込み、`http.response`型を作成するのが`readRequest`メソッドの仕事です。


```
// A response represents the server side of an HTTP response.
type response struct {
    conn    *conn
    req *Request // request for this response
    // (以下略)
}
```


出典:net/http/server.go

### 4-2. `http.serverHandler.ServeHTTP`メソッドの呼び出し

リクエスト内容を得ることができたら、いよいよハンドリングに入っていきます。\
`http.conn`のフィールドに含まれていた`http.Server`を、`http.serverHandler`型インスタンスにラップした上で`ServeHTTP`メソッドを呼び出します。


```
type serverHandler struct {
    srv *Server
}

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request)
```


出典:net/http/server.go

また、`ServeHTTP`メソッド呼び出しの際に渡している引数に注目すると、先ほど入手した`http.response`型が使われていることも特筆に値するでしょう。


```
// 再掲
func (c *conn) serve(ctx context.Context) {
    // (一部抜粋)
    w, err := c.readRequest(ctx)
    serverHandler{c.server}.ServeHTTP(w, w.req)
}
```



`http.Server`型をわざわざ`http.serverHandler`型にキャストすることによって、`http.Handler`インターフェースを満たすようになります。


## 5. `http.serverHandler`型の`ServerHTTP`メソッド

それでは、`http.serverHandler.ServerHTTP`メソッドの中身を見ていきましょう。


```
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
    // 一部抜粋
    handler := sh.srv.Handler
    handler.ServeHTTP(rw, req)
}
```


出典:net/http/server.go

1.  `sh.srv.Handler`で、`http.Handler`型インターフェースを得る
2.  Handlerインターフェースのメソッド、`ServeHTTP`を呼ぶ

![](https://storage.googleapis.com/zenn-user-upload/657ef5708bea0ee8a499328a.png)

### 5-1. `sh.srv.Handler`の取り出し

`http.serverHandler`の中には`http.Server`が存在し、そして`http.Server`の中には`http.Handler`が存在します。\
このハンドラを明示的に取り出しています。


```
type serverHandler struct {
    srv *Server
}
type Server struct {
    Handler Handler // これを取り出している
    // (以下略)
}
```


### 5-2. `http.Handler.ServeHTTP`メソッドの実行

`http.Handler`型というのは`ServeHTTP`メソッドを持つインターフェースです。\
上で取り出した`http.Handler`に対して、このメソッドを呼び出しています。


```
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```


出典:pkg.go.dev -
net/http#Handler

しかし、このインターフェースを満たす具体型は一体何なのでしょうか。

実は今までのコードをよくよく見返してみると、`sh.srv.Handler`で得られた`http.Handler`は、`http.ListenAndServe`関数を呼んだときの第二引数であるということがわかります。\
![](https://storage.googleapis.com/zenn-user-upload/2d1ff084770ebac5ebe2a073.png)
そのため、もしも


```
http.ListenAndServe(":8080", nil)
```


このようにサーバーを起動していた場合には、ここでの`http.Handler`を満たす具体型は、パッケージ変数`DefaultServeMux`の`http.ServeMux`型となります。

# 次章予告

ここまではサーバーの起動作業、具体的には

-   `net.Conn`を入手して、リクエストを受け取る体制を整える
-   ハンドラ関数の第二引数に渡す`http.response`型の用意
-   `http.ListenAndServe`関数の第二引数(今回は`nil`であるため`DefaultServeMux`となる)で渡されたルーティングハンドラの起動

までを追っていきました。

次章では、この続きを追いやすくするために、ルーティングハンドラである`DefaultServeMux`そのものについて詳しく掘り下げていきます。




# デフォルトでのルーティング処理の詳細

# この章について

前章にて「サーバー起動時に`http.ListenAndServe(":8080", nil)`とした場合、ルーティングハンドラはデフォルトで`net/http`パッケージ変数`DefaultServeMux`が使われる」という話をしました。

ここでは、この`DefaultServeMux`は何者なのかについて詳しく説明したいと思いいます。

# `DefaultServeMux`の定義・役割

## 定義

`DefaultServeMux`は、`net/http`パッケージ内に存在する公開グローバル変数です。


```
// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux
```


出典:net/http/server.go

`ServeMux`型の型定義は以下のようになっています。


```
type ServeMux struct {
    mu    sync.RWMutex
    m     map[string]muxEntry
    es    []muxEntry // slice of entries sorted from longest to shortest.
    hosts bool       // whether any patterns contain hostnames
}

type muxEntry struct {
    h       Handler
    pattern string
}
```


出典:net/http/server.go

## 役割

定義だけ見ても、`DefaultServeMux`で何を実現しているのかわかりにくいと思います。

実は`DefaultServeMux`の役割は、`ServeMux`の`m`フィールドが中心部分です。\
`m`フィールドの`map`には、「URLパスー開発者が`http.HandleFunc`関数で登録したハンドラ関数」の対応関係が格納されています。

Goのhttpサーバーは、`http.ListenAndServe`の第二引数`nil`の場合では`DefaultServeMux`内に格納された情報を使って、ルーティングを行っているのです。

# ハンドラの登録

まずは、`DefaultServeMux`に開発者が書いたハンドラが登録されるまでの流れを追ってみましょう。

開発者が書いた`func(w http.ResponseWriter, _ *http.Request)`という形のハンドラを登録するには、`http.HandleFunc`関数に対応するURLパスと一緒に渡してやることになります。


```
func main() {
    h1 := func(w http.ResponseWriter, _ *http.Request) {
        io.WriteString(w, "Hello from a HandleFunc #1!\n")
    }

    http.HandleFunc("/", h1) // パス"/"に、ハンドラh1が対応

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```


## 1. `http.HandleFunc`関数

それでは、`http.HandleFunc`関数の中身を見てみましょう。


```
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}
```


出典:net/http/server.go

内部では、`DefaultServeMux`(`http.ServeMux`型)の`HandleFunc`メソッドを呼び出しているだけです。

![](https://storage.googleapis.com/zenn-user-upload/e5dd85aa84fa47524a749ca5.png)

## 2. `http.ServeMux.HandleFunc`メソッド

それでは、`http.ServeMux.HandleFunc`メソッドの中身を見てみましょう。


```
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    if handler == nil {
        panic("http: nil handler")
    }
    mux.Handle(pattern, HandlerFunc(handler))
}
```


出典:net/http/server.go

内部で行っているのは主に2つです。

1.  `func(ResponseWriter, *Request)`型を、`http.HandlerFunc`型にキャスト
2.  ↑で作った`http.HandlerFunc`型を引数にして、`http.ServeMux.Handle`メソッドを呼ぶ

![](https://storage.googleapis.com/zenn-user-upload/5dd64341d0428087d5b53b69.png)

## 3. `http.ServeMux.Handle`メソッド

それでは、`http.ServeMux.Handle`メソッドの中を今度は覗いてみましょう。


```
func (mux *ServeMux) Handle(pattern string, handler Handler) {
    // (一部抜粋)
    e := muxEntry{h: handler, pattern: pattern}
    mux.m[pattern] = e
}
```


出典:net/http/server.go

ここで、`DefaultServeMux`の`m`フィールドに「URLパスー開発者が`http.HandleFunc`関数で登録したハンドラ関数」の対応関係を登録しています。

![](https://storage.googleapis.com/zenn-user-upload/c48843293a8909272404b115.png)

# `DefaultServeMux`によるルーティング

ここからは`DefaultServeMux`から、先ほど内部に登録したハンドラを探し当てるまでの処理を辿ってみましょう。

## 1. `http.ServeMux`の`ServeHTTP`メソッド

`DefaultServeMux`を使用したルーティング依頼は、`ServeHTTP`メソッドで行われます。

![](https://storage.googleapis.com/zenn-user-upload/0d73ed8e9a402db5a7dbe4ee.png)

このことは、前章の終わりが`http.Handler`インターフェースの`ServeHTTP`メソッドだったことを思い出してもらえると、このことが理解できると思います。\
`http.ServeMux`型は`ServeHTTP`メソッドを持つので、`http.Handler`インターフェースを満たします。

それでは、`http.ServeMux.ServeHTTP`メソッドの中身を見てみましょう。


```
// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
    // 一部抜粋
    h, _ := mux.Handler(r)
    h.ServeHTTP(w, r)
}
```


出典:net/http/server.go

ここで行っているのは次の2つです。

1.  `mux.Handler`メソッドで、リクエストにあったハンドラ(`http.Handler`インターフェース)を取り出す
2.  ↑で取り出したハンドラの`ServeHTTP`メソッドを呼び出す

![](https://storage.googleapis.com/zenn-user-upload/74d9b40979233c056b89aae1.png)

### 1-1. `http.ServeMux`の`Handler`メソッド

`mux.Handler`メソッド内では、どのようにリクエストに沿ったハンドラを取り出しているのでしょうか。\
それを知るために、`http.ServeMux.Handler`の中身を見てみましょう。


```
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
    // 一部抜粋
    return mux.handler(host, r.URL.Path)
}
```


出典:net/http/server.go

最終的に非公開メソッド`handler`メソッドが呼ばれています。

![](https://storage.googleapis.com/zenn-user-upload/530cbc84955e4f0733ce6db0.png)

### 1-2. `http.ServeMux`の`handler`メソッド

`http.ServeMux.handler`の中身は、以下のようになっています。


```
// handler is the main implementation of Handler.
// The path is known to be in canonical form, except for CONNECT methods.
func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
    // 一部抜粋
    if mux.hosts {
        h, pattern = mux.match(host + path)
    }
    if h == nil {
        h, pattern = mux.match(path)
    }
    if h == nil {
        h, pattern = NotFoundHandler(), ""
    }
    return
}
```


出典:net/http/server.go

`http.ServeMux.match`メソッドから得られるハンドラが返り値になっていることが確認できます。

![](https://storage.googleapis.com/zenn-user-upload/b088b99cf96bcd1f40cc9ca9.png)

### 1-3. `http.ServeMux`の`match`メソッド

そしてこの`http.ServeMux.match`メソッドが、「URLパス→ハンドラ」の対応検索を`DefaultServeMux`の`m`フィールドを使って行っている部分です。


```
// Find a handler on a handler map given a path string.
// Most-specific (longest) pattern wins.
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
    // Check for exact match first.
    v, ok := mux.m[path]
    if ok {
        return v.h, v.pattern
    }

    // Check for longest valid match.  mux.es contains all patterns
    // that end in / sorted from longest to shortest.
    for _, e := range mux.es {
        if strings.HasPrefix(path, e.pattern) {
            return e.h, e.pattern
        }
    }
    return nil, ""
}
```


出典:net/http/server.go

## 2. `http.Handler.ServeHTTP`メソッドの実行

`http.ServeMux.match`関数から得られた、ユーザーが登録したハンドラ関数(`http.Handler`インターフェース型)は、最終的には自身の`ServeHTTP`メソッドによって実行されることになります。


```
// 再掲
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
    // 一部抜粋
    h, _ := mux.Handler(r) // mux.match関数によってハンドラを探す
    h.ServeHTTP(w, r) // 実行
}
```


# まとめ

ルーティングハンドラである`DefaultServeMux`と、ユーザーが登録したハンドラ関数の対応関係は、以下のようにまとめられます。\
![](https://storage.googleapis.com/zenn-user-upload/533c8bb9af26d46da8e2eea4.png)

# 次章予告

次章では、「ルーティングハンドラによって取り出されたユーザー登録ハンドラ内で、どのようにレスポンスを返す処理を行っているのか」について掘り下げていきます。




# ハンドラによるレスポンス返却の詳細

# この章について

前2章を使って、

-   httpサーバーを起動し、
-   `DefaultServeMux`を使って、リクエストを適切なハンドラにルーティングする

ところまで追ってきました。

この章では、ルーティング後の話「ハンドラ内でどのようにしてレスポンスを作成し、返しているのか」について説明します。

# ハンドラ関数のおさらい

おさらいとして、ユーザーがサーバーに登録するハンドラの形をもう一度見てみます。


```
func main() {
    h1 := func(w http.ResponseWriter, _ *http.Request) {
        io.WriteString(w, "Hello from a HandleFunc #1!\n")
    }

    http.HandleFunc("/", h1) // パス"/"に、ハンドラh1が対応

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```


ハンドラ`h1`は、`func(w http.ResponseWriter, _ *http.Request)`というシグネチャをしています。

第二引数は、ハンドラが処理するリクエストが、`http.Request`型の形で入っているのだろうなと容易に想像がつきます。\
そのため、ここでは第一引数である`http.RewponseWriter`に注目します。

# 第一引数 - `http.ResponseWriter`

## 定義


```
type ResponseWriter interface {
    Header() Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int)
}
```


出典:pkg.go.dev -
net/http#ResponseWriter

`http.RewponseWriter`は、上記3つのメソッドを持つインターフェース型として定義されています。

ここで一つ疑問が生じます。\
ハンドラが受け取る`http.RewponseWriter`型第一引数の、実体型は何になるのでしょうか。

これはインターフェースです。これを満たす実体は何でしょうか。

## `http.ResponseWriter`インターフェースの実体型

`http.ResponseWriter`インターフェースの実体型を探すためには、`http.ListenAndServe`関数を呼んでから、この個別ハンドラの`ServeHTTP`メソッドが呼ばれるまでの変数の流れを順に追っていく必要があります。

以下の図は、それをまとめたものです。\
![](https://storage.googleapis.com/zenn-user-upload/deaebf46c7575b36c774a3a1.png)

ここから、図の下部にある`http.ResponseWriter`の大元は、2章前の`readRequest`メソッドにて登場した`http.response`型だということがわかります。

## `http.response`型

この`http.response`型の中には、サーバー起動の際に取得した`net.Conn`が含まれています。


```
// A response represents the server side of an HTTP response.
type response struct {
    conn    *conn
    req *Request // request for this response
    // (以下略)
}

// A conn represents the server side of an HTTP connection.
type conn struct {
    server *Server
    rwc net.Conn
    // (以下略)
}
```


そのため、`http.response.Write()`メソッドを呼ばれたときに実行されるのは、現在のコネクションである`net.Conn`の`Write`メソッドとなります。

したがって、


```
h1 := func(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "Hello from a HandleFunc #1!\n")
}
```


のように`http.ResponseWriter`に書き込まれた内容は、ネットワークを通じて返却するレスポンスへの書き込みとなるわけです。

## (210919追記)

Hiroaki
Nakamura(\@hnakamur2)さんから、「`http.response.Write()`メソッドを呼んだ後にネットワーク書き込みにたどり着くまで」についての補足情報をいただきました。

1.  非公開の`http.response.write()`メソッドが呼ばれる
2.  その中で、`http.response`型内部にある`bufio.Writer`の`Write`メソッドが呼ばれる
3.  `http.response`型内部の`bufio.Writer`インターフェースの具体型は、本記事3章で`http.response`型を取得するときに呼んだ`http.conn.readRequest`メソッドにて、`http.response.cw`フィールド(`http.chunkWriter`型)がセットされている
4.  `http.chunkWriter`型の`Write`メソッドにてネットワーク書き込みが行われ、この`Write`メソッドの中身を掘っていくと`net.Conn.Write`メソッドにたどり着く

ということです。情報ありがとうございました。




# GoでのWebサーバー起動の全体図

# この章について

前章までは、実際にコード内で呼ばれている関数・メソッドを網羅する形で処理の流れを追っていきました。\
そこで作った図は「正確」ではあるのですが、インターフェースや具体型が入り混じっており、その分大枠は掴みづらいものになっています。

そのためここでは、上で紹介した2つのフェーズの重要ポイントだけに絞る形で、処理の大枠をまとめ直してみます。

# 2つのフェーズ

GoでWebサーバーを起動させるときの処理は、大きく2つのフェーズに分けることができます。

1.  `http.Server`型や`net.Conn`インターフェースの作成といった、サーバーの起動処理
2.  実際に受信したリクエストをハンドラに処理させる、リクエストハンドリング

![](https://storage.googleapis.com/zenn-user-upload/968a1c46a20e4a9b0915603c.png)

# 処理の大枠

ここでは、上で紹介した2つのフェーズの大枠を述べていきます。

## 「インターフェース」で見る

処理の重要ポイントだけ抽出するには、メソッドセットの形である程度の抽象化がなされているインターフェースに着目するのがいいです。\
すると、処理の大枠は下図のようにまとめることができます。\
![](https://storage.googleapis.com/zenn-user-upload/8709d10f92c067ba7604793d.png)

### 1. サーバー起動

サーバーの起動部分で、最初に呼び出されるハンドラを内部に持つ`http.Server`型と、http通信をするための`net.Conn`インターフェースを作成しています。\
`net.Conn`が`http.Server`型の外にあるのは、おそらく依存性注入の観点での設計です。

-   `http.Server`型が持つルーティング情報は、どの環境で動かしたとしても不変なもの
-   `net.Conn`が持つURLホストやポートといったネットワーク環境情報は、状況によって変わる

これを踏まえて、もしURLやコネクションが変わったとしても`http.Server`型を作り直さなくてもいいようにしているのです。

### 2. リクエストハンドリング

実際にリクエストを受けて、レスポンスを返す段階になると、`http.Server`型は`ServeHTTP`メソッドがある`http.serverHandler`型にキャストされた上で、その`ServeHTTP`メソッドを呼び出すことでリクエストを捌いていきます。\
`serverHandler`型から最初に呼び出される`http.Handler`は、`http.ListenAndServe`の第二引数に渡されたルーティングハンドラです(=デフォルトだと`DefaultServeMux`)。

リクエストを受け取った`http.Handler`は、リクエストパスを見て、他の`http.Handler`に処理を委譲するか、自身でレスポンス作成をするかのどちらかの処理を行います。

## 具体型で見る

インターフェースで見た場合、リクエストをハンドルする部品は全て`http.Handler`でした。\
「他の`http.Handler`に処理を移譲するハンドラ」と「自身でリクエストを処理するハンドラ」の違いは一体なんなのでしょうか。

それをわかりやすくするために、上記の図を`http.Handler`インターフェースを満たしうる実体型で書き換えました。\
![](https://storage.googleapis.com/zenn-user-upload/7833d407b8e0eecb7ee04a24.png)

`http.Handler`インターフェース部分の具体型として使われているのは、大きく分けて二種類です。

-   `http.ServeMux`型:
    ルーティングハンドラ。リクエストパスをみて、他のハンドラに処理を振り分ける役割を担う。
-   `http.HandlerFunc`型:
    ユーザーが書いたhttpハンドラ。実際にレスポンス内容を作成し、`net.Conn`に書き込む役割を担う。


`http.serverHandler`型も`http.Handler`インターフェースを満たす型であるので、

-   処理の起点となる初めの`http.serverHandler`から別の`http.serverHandler`にハンドリング
-   `http.ServeMux`型から`http.serverHandler`にハンドリング

ということも理論上は可能です。\
ただし「あるサーバーから別のサーバーにハンドリング」というユースケースが現実的にありうるかどうかは疑問です(少なくとも筆者は思いつきません)。


「`http.ServeMux`型にするか、`http.HandlerFunc`型にするか」の選択イメージについては、以下の図のように「パス`/users`以降は別のハンドラに任せる」というようなハンドリングをする場合を思い浮かべてもらえればわかりやすいかと思います。

![](https://storage.googleapis.com/zenn-user-upload/37a28c299d663f2f1e71f8cc.png)




# おわりに

# おわりに

というわけで、標準`net/http`パッケージに絞ったWebサーバーの仕組みを掘り下げてきましたが、いかがでしたでしょうか。

`main`関数内で`http.ListenAndServe()`と書くだけで簡単にサーバーが起動できる裏では、

-   リクエストを受け取る・レスポンスを返すためのコネクションインターフェース(`net.Conn`)を生成
-   そこからリクエストを受け取ったら、`http.Handler`を乗り継いでレスポンスを処理するハンドラまで、リクエスト情報とネットワークインターフェースを伝達
-   ハンドラ中で作成したレスポンスを、`net.Conn`に書き込み

まで行う処理をうまく行ってくれています。

そこに至るまでの間も、`http.ResponseWriter`や`http.Handler`といったインターフェースを多用する柔軟な実装を行っており、

-   サーバーを動かすネットワーク環境(ホストURL、ポート)
-   ルーティング構成(直接ハンドラをぶら下げるか、内部にもルーターをいくつか繋げる形にするか)

といった多種多様なサーバーに適用できるように`net/http`は作られているのです。

また、今回は触れませんでしたが、Webサーバーの実装を具体型に頼るのではなく、`http.Handler`インターフェースを多用する形にしたことによって、例えば`gorilla/mux`といった外部ルーティングライブラリを導入したとしても、


```
import (
    "github.com/gorilla/mux"
)

func main() {
    // ハンドラh1, h2を用意(略)

    r := mux.NewRouter()  // 明示的にgorilla/muxのルータを用意し、
    r.HandleFunc("/", h1) // そのルータのメソッドを呼んで登録し
    r.HandleFunc("/endpoint", h2)

    log.Fatal(http.ListenAndServe(":8080", r)) // それを最初のルータとして使用
}
```


このようにユーザー側が大きくコードを変えることなく、自分たちが使いたいルーティングシステムと`net/http`の仕組みを共存させて使うことができるようになります。\
ユーザー側がどんなライブラリの部品を渡してきたとしても`net/http`パッケージ側がそれに対応できるあたりでも、GoのWebサーバーの柔軟さを感じていただけると思います。

この記事を通して、GoでWebサーバーを起動させた裏側の話と、その設計の柔軟さしなやかさについて少しでも「わかった！」となってもらえたら嬉しいです。

コメントによる編集リクエスト・情報提供等大歓迎です。\
連絡先: 作者Twitter \@saki_engineer

# 後日談

このZenn本のメイキングについて、こちらの記事で公開しました。\






<https://zenn.dev/hsaki/articles/go-serverinternal-making>

# 参考/関連文献

今回の話を書くにあたって参考にした文献と、本記事では触れなかった`net/http`以外のWeb周り関連サードパーティパッケージについて軽く紹介したいと思います。

## 公式ドキュメント

### `net/http`






<https://pkg.go.dev/net/http>

`net/http`について深く掘り下げたいなら、とにもかくにも公式ドキュメントをあたりましょう。\
Webサーバー周りは需要が多い分野であるため、サンプルコードも豊富に掲載されています。

## 一般のブログ

### Future Tech Blog - Goのおすすめのフレームワークは`net/http`






<https://future-architect.github.io/articles/20210714a/>

> 僕としてはGoのおすすめのフレームワークを聞かれたら、標準ライブラリのnet/httpと答えるようにしています。\
> というよりも、Goの他のフレームワークと呼ばれているものは、このnet/httpのラッパーでしかないからです。

最初の2文に言いたいことが全て詰まってますね。\
メッセージ性の強さゆえに一時期ものすごく話題になった記事です。この記事の主張については私も賛成で、「まあ流石にパスパラメータが絡んだら`gorilla/mux`はいれるけど、フルスタックフレームワークは入れないなあ」派です。

また、`net/http`内で`http.Handler`インターフェースを多用するとどう設計の柔軟さが生まれているのかについて、非常に直感的でわかりやすい絵を用いて説明しているのも必見ポイントです。

### tenntenn.dev - GoでオススメのWebフレームワークを聞かれることが多い






<https://tenntenn.dev/ja/posts/2021-06-27-webframework/>

外部モジュール導入による複雑さの増大は、Goのシンプルさを阻害するゆえによく考えた方がいい、`net/http`から始めてもいいのでは？という記事です。\
書き手は違えど、主張は先ほどの記事と似ていますね(公開時期はこちらの方が少しだけ前ですが)。

### なぜGo言語のHTTPサーバーでクライアント切断が検知できるのか調べた






<https://zenn.dev/najeira/articles/2020-12-17-zenn-chiesh7noijeequii7ae>

今回はわかりやすさのために正常系に絞って説明しましたが、この記事では「contextによるキャンセル処理」がどうnet/httpの中に組み込まれているのかという触りの部分が書かれています。

## サードパーティパッケージ

### ルーティングライブラリ

`net/http`の標準ルータ`DefaultServeMux`を置き換えるような使い方をするライブラリです。


```
// 第二引数を置き換える
-log.Fatal(http.ListenAndServe(":8080", nil))
+log.Fatal(http.ListenAndServe(":8080", pkg.router))
```


`DefaultServeMux`には難しい、パスパラメータの抽出などが簡単に行うことができます。

主に以下の2つがよく聞くものになるでしょうか。

#### `gorilla/mux`






<https://www.gorillatoolkit.org/>\






<https://github.com/gorilla/mux>

#### `go-chi/chi`






<https://github.com/go-chi/chi>

### Webフレームワーク

ルーティングライブラリが、サーバー起動・ハンドラ登録の部分に関しては`net/http`の仕組みをそのまま使うのに対し、Webフレームワークになるとその部分までもパッケージ独自のやり方で行うようになります。


```
// echoの場合
import (  
    "github.com/labstack/echo/v4"
)

func main() {
    // net/httpの痕跡は表面上は見られない
    e := echo.New()
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":1323"))
}
```


以下の2つはよく聞きます。\
が、特徴に関しては詳しくないので筆者にはわかりません。

#### `labstack/echo`






<https://echo.labstack.com/>\






<https://github.com/labstack/echo>

#### `gin-gonic/gin`






<https://gin-gonic.com/>\






<https://github.com/gin-gonic/gin>

## LTスライド

### HTTPルーティングライブラリ入門






<https://speakerdeck.com/hikaru7719/http-routing-library>

golang.tokyo#31にて行われたセッションです。\
`net/http`,`gorilla/mux`,`chi`の3つのルーティングライブラリについて、それぞれを使用したコード全体の比較した上で、さらに

-   ルーティングアルゴリズム
-   パスパラメータの取得方法

の入門のような内容が書かれています。




