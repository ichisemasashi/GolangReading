# Goから学ぶI/O


GoにはI/Oに関わるパッケージが数多く存在します。io, os, bufio,
fmtなどなど......。これらの立ち位置や、I/O実行の裏で何が起こっているのか本当に理解していますか？この本では、この問への答えをまとめました。

# はじめに

# この本について

この本では、Go言語で扱えるI/Oについてまとめています。

Go言語I/Oを扱うためのパッケージとしては、ドンピシャのものとしては`io`パッケージがあります。\
しかし、例えば実際にファイルを読み書きしようとするときに使うのは、`os`パッケージの`os.File`型まわりのメソッドです。\
標準入力・出力を扱おうとすると`fmt`パッケージが手っ取り早いですし、また速さを求める場面では`bufio`パッケージのスキャナを使うということもあるでしょう。\
このように、「I/O」といってもGoでそれに関わるパッケージは幅広いのが現状です。

また、ファイルオブジェクト`f`に対して`f.Read()`とかいう「おまじない」と唱えるだけで、なんでファイルの中身が取得できるの？一体裏で何が起こっているの？という疑問を感じている方もいるかと思います。

ここでは

-   `os`や`io`とかいっぱいあるけど、それぞれどういう関係なの？
-   標準入力・出力を扱うときに`fmt`と`bufio`はどっちがいいの？
-   そもそも`bufio`パッケージって何者？
-   GoでやったI/Oまわりの操作は、実現のために裏で何が起こっているの？\
    こういったことを一から解説していきます。

## 使用する環境・バージョン

-   OS: macOS Mojave 10.14.5
-   go version go1.16.2 darwin/amd64

## 読者に要求する前提知識

-   Goの基本的な文法の読み書きができること
-   基本情報技術者試験くらいのIT前提知識




# ファイルの読み書き

# はじめに

Goでファイルの読み書きに関する処理は`os`パッケージ中に存在する`File`型のメソッドで行います。

この章では

-   `os.File`型って一体何者？
-   読み書きってどうやってするの？
-   低レイヤでは何が起こっているの？\
    ということについてまとめていきます。

# ファイルオブジェクト

`os`パッケージには`os.File`型が存在し、Goでファイルを扱うときはこれが元となります。


``` go
type File struct {
    *file // os specific
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/os/types.go#16>

`os.File`型の実際の実装は`os.file`型という非公開型で行われており、その内部構造については外から直接見ることができないようになっています。


このように「公開する構造体の中身を隠したい場合に、隠す中身を非公開の構造体型にしてまとめて、公開型構造体に埋め込む」という手段はGoの標準パッケージ内ではよく見られる手法です。


# ファイルを開く(open)

## 読み込み権限onlyで開く

Go言語でファイルを扱い読み書きするためには、まずはそのファイルを\"open\"して、`os.File`型を取得しなくてはいけません。

`os.File`型を得るためには、`os.Open(ファイルパス)`関数を使います。


``` go
f, err := os.Open("text.txt")
```


得られる第一返り値`f`が`os.File`型のファイルオブジェクトです。

`os.Open()`関数について、ドキュメントでは以下のように書かれています。

> `Open` opens the named file for reading. If successful, methods on the
> returned file can be used for reading; the associated file descriptor
> has mode O_RDONLY.
>
> (訳)
> `Open`関数は、名前付きのファイルを読み込み専用で開きます。`Open`が成功すれば、返り値として得たファイルオブジェクトのメソッドを中身の読み込みのために使うことができます。`Open`から得たファイルは、Linuxでいう`O_RDONLY`フラグがついた状態になっています。
>
> 出典:pkg.go.dev - os package

## 書き込み権限付きで開く

書き込み権限がついた状態のファイルが欲しい場合、`os.Create(ファイルパス)`関数を使います。


``` go
f, err := os.Create("write.txt")
```


`Open()`と同様に、これも第一返り値`f`が`os.File`型のファイルオブジェクトです。

\"create\"の名前を見ると「ファイルがない状態からの新規作成にしか対応してないのか？」と思う方もいるでしょうが、引数のファイルパスには既に存在しているファイルの名前も指定することができます。今回の場合、`write.txt`が既に存在してもしなくても、上のコードは正しく動作します。

ドキュメントに記載されている`os.Create()`の説明は以下のようになっています。

> Create creates or truncates the named file. If the file already
> exists, it is truncated. If the file does not exist, it is created
> with mode 0666 (before umask). If successful, methods on the returned
> File can be used for I/O; the associated file descriptor has mode
> O_RDWR.
>
> (訳)`Create()`関数は、名前付きファイルを作成するか、中身を空にして開きます。引数として指定したファイルが既に存在している場合、中身を空にして開くほうの動作がなされます。ファイルが存在していなかった場合は、`umask 0666`のパーミッションでファイルを作成します。`Create()`が成功すれば、返り値として得たファイルオブジェクトのメソッドをI/Oのために使うことができます。`Create`から得たファイルは、Linuxでいう`O_RDWR`フラグがついた状態になっています。
>
> 出典:pkg.go.dev - os package


truncateは、直訳が「切り捨てる」という動詞です。Linuxの文脈では、truncateは「ファイルサイズを指定したサイズにする」という意味で使われることが多いです。これには、ファイルサイズを大きくすることも小さくすることも含まれ、例えば10byteのファイルを20byteにする処理も、訳語に反しますが\"truncate\"です。ファイルサイズが指定されなかった場合、ファイルサイズ0にtruncateされるととられ、今回の`Create`の場合はこちらの動作になります。


# ファイル内容の読み込み(read)

同じディレクトリ中にある`text.txt`の内容をすべて読み込むという操作を考えます。


    Hello, world!
    Hello, Golang!


これをGoで行う場合、`os.File`型の`Read`メソッドを用いて以下のように実装できます。


``` go
// os.FileオブジェクトをOpen関数か何かで事前に得ておくとする
// 変数fがファイルオブジェクトとする

data := make([]byte, 1024)
count, err := f.Read(data)
if err != nil {
    fmt.Println(err)
    fmt.Println("fail to read file")
}

/*
--------------------------------
挙動の確認
--------------------------------
*/

fmt.Printf("read %d bytes:\n", count)
fmt.Println(string(data[:count]))

/*
出力結果

read 28 bytes:
Hello, world!
Hello, Golang!
*/
```


`Read(b []byte)`メソッドの引数としてとる`[]byte`スライスの中に、読み込まれたファイルの内容が格納されます。

また、`Read()`メソッドの第一返り値(上での`count`変数に値が格納)には、「`Read()`メソッドを実行した結果、何byteが読み込まれたか」が`int`型で入っています。\
そのため、`string(data[:count])`とすることで、ファイルから読み込まれた文字列をそのまま得ることができます。


`fmt.Println(string(data[:count]))`\
↓\
`fmt.Println(data[:count])`\
のようにprintする内容を変更すると、「文字列」ではなくて「文字列にエンコードする前のバイト列そのまま」が得られるので注意。\
(例)\
文字列→`"Hello, world!\nHello, Golang!"`\
エンコード前→バイト列`[72 101 108 108 111 44 32 119 111 114 108 100 33 10 72 101 108 108 111 44 32 71 111 108 97 110 103 33]`


# ファイルへの書き込み(write)

ファイルに何かを書き込むときは、`os.File`型の`Write()`メソッドを利用します。

実際に`write.txt`というテキストファイルに文字列を書き込むコードを実装してみます。


``` go
// fはos.Create()で得たファイルオブジェクトとします。

str := "write this file by Golang!"
data := []byte(str)
count, err := f.Write(data)
if err != nil {
    fmt.Println(err)
    fmt.Println("fail to write file")
}

/*
--------------------------------
挙動の確認
--------------------------------
*/
fmt.Printf("write %d bytes\n", count)
/*
出力結果
write 26 bytes
*/
```



    write this file by Golang!


`Write`メソッドの引数としてとる`[]byte`スライス(ここでは変数`data`)に格納されている内容が、ファイルにそのまま書き込まれることになります。\
ここでは引数に「文字列`write this file by Golang!`をバイト列にキャストしたもの」を使っているので、この文字列がそのまま`write.txt`に書き込まれます。

また、`Write`メソッドの第一返り値には、「メソッド実行の結果ファイルに何byte書き込まれたか」が`int`型で得られます。


`Write`メソッドを使う予定のファイルオブジェクトは、書き込み権限がついた`os.Create()`から作ったものでなくてはなりません。\
`os.Open()`で開いたファイルは読み込み専用なので、これに`Write`メソッドを使うと、以下のようなエラーが出ます。\
`write write.txt: bad file descriptor`


# ファイルを閉じる(close)

## 基本

ファイルを閉じるときは`os.File`型`Close`メソッドを用います。


``` go
f, err := os.Open("text.txt")
if err != nil {
    fmt.Println("cannot open the file")
}
defer f.Close()

// 以下read処理等を書く
```


上のコードでは、`Close()`メソッドは`defer`を使って呼んでいます。\
一般的に、ファイルというのは「開いて使わなくなったら必ず閉じるもの」なので、`Close()`は`defer`での呼び出し予約と非常に相性がいいメソッドです。

## 応用

ところで、`Close`メソッドの定義をドキュメントで見てみると、以下のようになっています。


``` go
func (f *File) Close() error
```


出典:pkg.go.dev - os#File.Close
このように、実は返り値にエラーがあるのです。

ファイルを開いた後に行う操作が「読み込み」だけの場合、元のファイルはそのままですから`Close()`に失敗するということはほとんどありません。\
そのため、基本の節では`Close`メソッドから返ってくるエラーをさらっと握り潰してしまいました。

しかし、開いた後に行う操作が「書き込み」のような元のファイルに影響を与えるような操作だった場合、その処理が正常終了しないと`Close`できない、という状態に陥ることがあります。\
そのため、`Write`メソッドを使う場合は`Close`の返り値エラーをきちんと処理すべきです。

`defer`を使いつつエラーを扱うためには、以下のように無名関数を使います。


```
f, err := os.Create("write.txt")
if err != nil {
    fmt.Println("cannot open the file")
}
- defer f.Close()
+ defer func(){
+    err := f.Close()
+    if err != nil {
+        fmt.Println(err)
+    }
+ }

// 以下write処理等を書く
```


# 低レイヤで何が起きているのか

Goのコード上で`os.Open()`だったり`f.Read()`だったりを「おまじない」のように唱えることで、実際のファイルを扱うことができるのは一体どうなっているのでしょうか。\
これをよく知るためには、OSカーネルへと続く低レイヤなところに視点を下ろす必要があります。\
本章では`os`パッケージのコードを深く掘り下げることでこれを探っていきます。

## ファイルオブジェクト

`os.File`型の中身は以下のようになっています。


``` go
type file struct {
    pfd         poll.FD
    name        string
    dirinfo     *dirInfo // nil unless directory being read
    nonblock    bool     // whether we set nonblocking mode
    stdoutOrErr bool     // whether this is stdout or stderr
    appendMode  bool     // whether file is opened for appending
}
```


(`os.File`型の中身が、非公開の構造体`os.file`型であるのは前述した通りです)\
出典:<https://go.googlesource.com/go/+/go1.16.3/src/os/file_unix.go#54>

この中で重要なのは`pfd`\[1\]フィールドです。

Linuxカーネルプロセス内部では、openしたファイル1つに非負整数識別子1つを対応付けて管理しており、この非負整数のことをfd(ファイルディスクリプタ)と呼んでいます。\
`poll`パッケージの`FD`型はこのfdをGo言語上で具現化した構造体なのです。

> FD is a file descriptor. The net and os packages use this type as a
> field of a larger type representing a network connection or OS file.\
> (訳)`FD`型はファイルディスクリプタです。`net`や`os`パッケージでは、ネットワークコネクションやファイルを表す構造体の内部フィールドとしてこの型を使用しています。\
> 出典:pkg.go.dev - internal/poll
> package

`FD`型の定義は以下のようになっていて、この`Sysfd`という`int`型のフィールドがfdの数字そのものを表しています。


``` go
type FD struct {
    // System file descriptor. Immutable until Close.
    Sysfd int

    // Whether this is a streaming descriptor, as opposed to a
    // packet-based descriptor like a UDP socket. Immutable.
    IsStream bool

    // Whether a zero byte read indicates EOF. This is false for a
    // message based socket connection.
    ZeroReadIsEOF bool

    // contains filtered or unexported fields
}
```


出典:pkg.go.dev -
internal/poll#FD


ちなみにカーネルでは、openしていない全てのファイルに対しても整数の識別子をつけて管理しており、これをinode番号といいます。\
fdはそれとは区別された概念で、こちらは「プロセス中でopenしたファイルに対して順番に割り当てられる番号」です。

そのため、同じファイルを開いたらいつもfdが同じ番号になる、という代物ではありません。\
あるプログラムで`read.txt`を開いたらfdが3になったけど、別のときに別のプログラムで`read.txt`を開いたらfdが4になる、ということは普通に存在します。


## ファイルオープン

`os.Open()`実装の中身をこれからみていきます。

まず、`os.Open`自体は、同じ`os`パッケージの`OpenFile`関数を呼んでいるだけです。


``` go
func Open(name string) (*File, error) {
    return OpenFile(name, O_RDONLY, 0)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/os/file.go#310>


ちなみに`os.Create()`も内部で`OpenFile`を呼んでいます。ただし、ファイルに書き込み権限をつけるため、関数に渡している引数が違います。

というより`OpenFile`関数そのものが「ファイルを特定の権限で開く」ための一般的な操作を規定したもので、`os.Open`や`os.Create`はこれをユーザーがよく使う引数でラップしただけ、というのが本来の位置付けです。


`os.OpenFile`関数の中身を見ると、非公開関数`openFileNolog`を呼んでいるのがわかります。


``` go
func OpenFile(name string, flag int, perm FileMode) (*File, error) {
    // (略)
    f, err := openFileNolog(name, flag, perm)
    // (略)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/os/file.go#329>

この`openFileNoLog`関数をみると、内部では`syscall.Open()`という`syscall`パッケージの関数が呼ばれています。


``` go
func openFileNolog(name string, flag int, perm FileMode) (*File, error) {
    // (略)
    var r int
    for {
        var e error
        r, e = syscall.Open(name, flag|syscall.O_CLOEXEC, syscallMode(perm))
        if e == nil {
            break
        }
        // (略:EINTRエラーを握り潰す処理)
    }
    // (略)
    return newFile(uintptr(r), name, kindOpenFile), nil
}
```


出典:<https://go.googlesource.com/go/+/go1.16.2/src/os/file_unix.go#205>


EINTRは、処理中に割り込み信号(ユーザーによるCtrl+Cなど)があったというエラー番号のこと。


`openFileNolog`関数の返り値とするために、`syscall.Open`から得られた返り値`r`をfdとする`os.File`型を生成しています。\
言い換えると、「ファイルのfdを得る」という根本的な操作をしているのは`syscall.Open`関数です。

この`syscall`パッケージでは、OSカーネルへのシステムコールをGoのソースコードから呼び出すためのインターフェースを定義しています。

> Package syscall contains an interface to the low-level operating
> system primitives.\
> 出典:pkg.go.dev - syscall package

そしてこの`syscall.Open`関数は、OSの`open`システムコールを呼び出すためのラッパーなのです。後の処理はカーネルがやってくれます。

Linuxの場合、システムコール`open()`は、指定したパスのファイルを指定したアクセスモードで開き、返り値としてfdを返すものです。


``` language-c
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>


int open(const char *pathname, int flags);
```


この引数`flags`に入れられるフラグとして`O_RDONLY`や`O_RDWR`があり、これによってopenしたファイルが読み込み専用になったり、読み書き可能になったりします。

## Readメソッド

次に、`os.File`型の`Read`メソッドを掘り下げてみましょう。

先述した通り、`os.File`型の実体は非公開の`os.file`型です。\
そしてこの`os.file`型の`Read`メソッドは、非公開メソッド`read`メソッドを経由して、その構造体のフィールドの一つ`pfd`(`poll.FD`型)の`Read`メソッドを呼んでいます。


``` go
// os.file型の公開Readメソッドの中身
func (f *File) Read(b []byte) (n int, err error) {
    // (中略)
    n, e := f.read(b)  // ここで読み込み(非公開readメソッドを呼び出し)
    return n, f.wrapErr("read", e)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/os/file.go#113>


``` go
// os.file型の非公開readメソッドの中身
func (f *File) read(b []byte) (n int, err error) {
    n, err = f.pfd.Read(b) // ここで読み込み
    // (中略)
    return n, err
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/os/file_posix.go#30>

この`poll.FD`型の`Read()`メソッドの内部実装で、`ignoringEINTRIO(syscall.Read, fd.Sysfd, p)`というコードが存在します。\
ここで呼ばれている`syscall.Read`関数が、OSカーネルの`read`システムコールのラッパーです。ここでGoと低レイヤとつながるのです。\
出典:<https://go.googlesource.com/go/+/go1.16.2/src/internal/poll/fd_unix.go#162>

順番をまとめると、`os.File`型の`Read`メソッドは以下のような実装となっています。

1.  `os.file`型の`Read`メソッドを呼ぶ
2.  1の中で`os.file`型の`read`メソッドを呼ぶ
3.  2の中で`poll.FD`型の`Read`メソッドを呼ぶ
4.  3の中で`syscall.Read`メソッドを呼ぶ
5.  OSカーネルのシステムコールで読み込み処理

## Writeメソッド

`os.File`型の`Write()`メソッドのほうも`Read`メソッドと同様の流れで実装されています。

1.  `os.file`型の`Write`メソッドを呼ぶ
2.  1の中で`os.file`型の`write`メソッドを呼ぶ
3.  2の中で`poll.FD`型の`Write`メソッドを呼ぶ
4.  3の中で`syscall.Write`メソッドを呼ぶ
5.  OSカーネルのシステムコールで書き込み処理

![](https://storage.googleapis.com/zenn-user-upload/rnaugoc5gva6ra9kkxzexjjgno2o)

## (おまけ)ファイルクローズ

ここまで見てきたファイル操作の裏には、どれもシステムコールがありました。\
なので「ファイルの`Close()`メソッドも、裏ではclose()のシステムコールを呼んでいるんでしょ？」と推測する方もいるかもしれません。

しかし実は、`os.File`型の`Close()`メソッドを掘り下げても、closeシステムコールに繋がる`syscall.Close`は出てきません。\
これはなぜかというと、ファイルオープンの時点で「ファイルオープンしたプロセスが終了したら、自動的にファイルを閉じてください」という`O_CLOEXEC`フラグを立てているからなのです。


``` go
// (再掲)
func openFileNolog(name string, flag int, perm FileMode) (*File, error) {
    // (略)
    // 第二引数が「フラグ」
    r, e = syscall.Open(name, flag|syscall.O_CLOEXEC, syscallMode(perm))
    // (略)
}
```


そのため、`Close()`メソッドがやっているのは

-   エラー処理
-   対応する`os.File`型を使えなくする後始末

という側面が強いです。

# まとめ

ここまでは、ファイルの読み書きについて取り上げました。\
ただし、「I/O」というのはファイルだけのものではありません。

次章では、「ファイルではないI/O」について扱いたいと思います。



脚注


1.  
    pfdはおそらくpollパッケージのFD型の略です。
    





# ネットワーク

# はじめに

この章ではネットワークについて扱います。\
「ネットワークにI/Oがなんの関係があるの？」と思う方もいるかもしれませんが、「サーバーからデータを受け取る」「クライアントからデータを送る」というのは、言い換えると「コネクションからデータを読み取る・書き込む」ともいえるのです。

`net`パッケージのドキュメントには以下のように記載されています。

> Package net provides a portable interface for **network I/O**,
> including TCP/IP, UDP, domain name resolution, and Unix domain
> sockets.\
> (訳)`net`パッケージでは、TCP/IP, UDP, DNS,
> UNIXドメインソケットを含むネットワークI/Oのインターフェース(移植性あり)を提供します。\
> 出典:pkg.go.dev - net package

ネットワークをI/Oと捉える言葉が明示されているのがわかります。

ここからは、TCP通信で短い文字列を送る・受け取るためのGoのコードについて解説していきます。

# ネットワークコネクション

ネットワーク通信においては、「クライアント-サーバー」間を繋ぐコネクションが形成されます。\
このコネクションパイプをGoで扱うインターフェースが`net.Conn`インターフェースです。


``` go
type Conn interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
    LocalAddr() Addr
    RemoteAddr() Addr
    SetDeadline(t time.Time) error
    SetReadDeadline(t time.Time) error
    SetWriteDeadline(t time.Time) error
}
```


出典:pkg.go.dev - net#Conn

`net.Conn`インターフェースは8つのメソッドセットで構成されており、これを満たす構造体としては`net`パッケージの中だけでも`net.IPConn`,
`net.TCPConn`, `net.UDPConn`, `net.UnixConn`があります。

# コネクションを取得

## サーバー側から取得する

サーバー側から`net.Conn`インターフェースを取得するためには、以下のような手順を踏みます。

1.  `net.Listen(通信プロトコル, アドレス)`関数から`net.Listener`型の変数(`ln`)を得る
2.  `ln`の`Accept()`メソッドを実行する


``` go
ln, err := net.Listen("tcp", ":8080")
if err != nil {
    fmt.Println("cannot listen", err)
}
conn, err := ln.Accept()
if err != nil {
    fmt.Println("cannot accept", err)
}
```


`conn`が`net.Conn`インターフェースの変数で、今回の場合、その実体はTCP通信のために使う`net.TCPConn`型構造体です。\
![](https://storage.googleapis.com/zenn-user-upload/o27hayivyrxb3f1sice2v688pifa)

## クライアント側から取得する

クライアント側から`net.Conn`インターフェースを取得するためには、`net.Dial(通信プロトコル, アドレス)`関数を実行します。


``` go
conn, err := net.Dial("tcp", "localhost:8080")
if err != nil {
    fmt.Println("error: ", err)
}
```


この`conn`も実体は`net.TCPConn`型です。\
![](https://storage.googleapis.com/zenn-user-upload/qpd98ckv0y2foe9uc312hzx9y591)

# サーバー側からのデータ発信

サーバー側から、TCPコネクションを使って文字列`"Hello, net pkg!"`を一回送信する処理は、`net.TCPConn`の`Write`メソッドを利用して以下のように実装されます。


``` go
// コネクションを得る
ln, err := net.Listen("tcp", ":8080")
if err != nil {
    fmt.Println("cannot listen", err)
}
conn, err := ln.Accept()
if err != nil {
    fmt.Println("cannot accept", err)
}

// ここから送信

str := "Hello, net pkg!"
data := []byte(str)
_, err = conn.Write(data)
if err != nil {
    fmt.Println("cannot write", err)
}
```


`Write`メソッドの挙動は、`os.File`型の`Write`メソッドのものとそう変わりません。\
引数にとった`[]byte`列の内容をコネクションに書き込み、そして何byte書き込めたかの値が第一返り値になります。

# クライアント側がデータ受信

クライアントがTCPコネクションから、文字列データを受け取るコードを`net.TCPConn`の`Read`メソッドを利用して書きます。


``` go
// コネクションを得る
conn, err := net.Dial("tcp", "localhost:8080")
if err != nil {
    fmt.Println("error: ", err)
}

// ここから読み取り

data := make([]byte, 1024)
count, _ := conn.Read(data)
fmt.Println(string(data[:count]))

// 出力結果
// Hello, net pkg!
```


`Read`メソッドの挙動も`os.File`の`Read`メソッドと同じです。\
引数にとった`[]byte`列の中に、コネクションから読み取った内容を入れて、そして何byte読めたかの値が第一返り値になります。

# 低レイヤで何が起きているのか

ここからは、`os.File`型のときにやったのと同様のコードリーディングを行います。\
ネットワークまわりのI/Oの実装では、どのようなシステムコールにつながっているのでしょうか。低レイヤの話に深く潜り込んでいきます。

## ネットワークコネクション(net.TCPConnの正体)

`net.TCPConn`構造体の正体は、非公開の構造体`net.conn`型です。


``` go
type TCPConn struct {
    conn
}
```


出典:<https://go.googlesource.com/go/+/go1.16.2/src/net/tcpsock.go#85>

そしてこの`net.conn`型の中身は、`netFD`型構造体そのものです。


``` go
type conn struct {
    fd *netFD
}
```


出典:<https://go.googlesource.com/go/+/go1.16.2/src/net/net.go#170>

この`netFD`型は一体何なのでしょうか。これも定義を見てみましょう。


``` go
type netFD struct {
    pfd poll.FD
    // immutable until Close
    family      int
    sotype      int
    isConnected bool // handshake completed or use of association with peer
    net         string
    laddr       Addr
    raddr       Addr
}
```


出典:<https://go.googlesource.com/go/+/go1.16.2/src/net/fd_posix.go#17>

前章で出てきた`poll.FD`型の`pfd`フィールドがここでも登場しました。これは一体どういうことでしょうか。

実はLinuxの設計思想として **\"everything-is-a-file philosophy\"**
というものがあります。これは、キーボードからの入力も、プリンターへの出力も、ハードディスクやネットワークからのI/Oもありとあらゆるものを全て「**OSのファイルシステム上にあるファイルへのI/Oとして捉える**」という思想です。\
今回のようなネットワークからのデータ読み取り・書き込みも、OS内部的には通常のファイルI/Oと変わらないのです。そのため、ネットワークコネクションに対しても、通常ファイルと同様にfdが与えられるのです。\
![](https://storage.googleapis.com/zenn-user-upload/av0ff3br57ap2iygje9p6qf1hnzw)

## コネクションオープン

では、通信するネットワークに対応するfdはどのように決まるのでしょうか。\
また、コネクションに対応したfdが入った`net.Conn`(ここでは`net.TCPConn`型構造体)はどのようにして得られるのでしょうか。

これを理解するためには、

-   クライアント側で`net.Dial()`を実行
-   サーバー側で`net.Listen()`→`ln.Accept()`を実行

それぞれにおいて裏で何が起きているのか、コードを読んで深掘りしていきましょう。

### クライアント側からのコネクションオープン

まずは、クライアント側から`net.Conn`を得るために呼ぶ`net.Dial(通信プロトコル, アドレス)`の中身をみてみます。\
すると、今私たちが欲しい「コネクションに割り当てられたfdをもつ`net.TCPConn`」を作っているのは、実質`net.Dialer`型の`DialContext`メソッドであることがわかります。


``` go
func Dial(network, address string) (Conn, error) {
    var d Dialer
    return d.Dial(network, address)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/net/dial.go#317>


``` go
func (d *Dialer) Dial(network, address string) (Conn, error) {
    return d.DialContext(context.Background(), network, address) // net.TCPConnを作っているのはここ
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/net/dial.go#347>

`net.Dialer`型の`DialContext`メソッドは、「引数として渡されたプロトコル・URL・ポート番号に対応した`net.Conn`を作る」ためのメソッドです。

> DialContext connects to the address on the named network using the
> provided context.\
> 出典:pkg.go.dev -
> net#Dialer.DialContext

この`DialContext`メソッドでやっていることは中々複雑なのですが、核としては

1.  `syscall.Socket`経由でシステムコールsocket()を呼んで、URLやポート番号からfdをゲットする
2.  1で得たfdを`poll.FD`型にする
3.  2で得た`poll.FD`型の`fd`を使い`newTCPConn(fd)`を実行→これが`TCPConn`になる

という流れです。\
![](https://storage.googleapis.com/zenn-user-upload/36rx6fx5dw4qd3o6lriieeji4ur2)
結局のところ、システムコールsocket()を内部で呼んで得たfdを`TCPConn`型にラップしている、ということです。

### サーバー側からのコネクションオープン

サーバー側で`net.Listen()`→`ln.Accept()`という手順を踏んだ場合は何が起こっているのでしょうか。\
`net.Listen()`関数の実装を確認してみます。


``` go
func Listen(network, address string) (Listener, error) {
    var lc ListenConfig
    return lc.Listen(context.Background(), network, address)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/net/dial.go#704>

`net.ListenConfig`型の`Listen`メソッドを内部で呼んでいます。\
この`Listenメソッド`の中身も中々複雑ですが、核は

1.  `syscall.Socket`経由でシステムコールsocket()を呼んで、URLやポート番号からfdをゲットする
2.  1で得たfdを内部フィールドに含んだ`TCPListener`型を生成し、返り値にする

となっています。\
ここでも、コネクションに対応したfdを得るからくりはsocket()システムコールです。

ですがまだ`net.Listener`が得られただけで、実際に通信に使う`TCPConn`構造体がまだです。\
実は、この「リスナーからコネクションを得る」ためのメソッドが`Accept()`メソッドなのです。その中身をみてみます。


``` go
func (l *TCPListener) Accept() (Conn, error) {
    // (略)
    c, err := l.accept()
    // (略)
    return c, nil
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/net/tcpsock.go#257>

内部では非公開メソッド`accept()`を呼んでいました。その中身は以下のようになっています。


``` go
func (ln *TCPListener) accept() (*TCPConn, error) {
    // リスナー本体からfdを取得
    fd, err := ln.fd.accept()
    // (略)

    // fdからTCPConnを作成
    tc := newTCPConn(fd)
    // (略)

    return tc, nil
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/net/tcpsock_posix.go#138>

要するに、「リスナーからコネクションを得る」=「リスナーからfdを取り出して、それを`TCPConn`にラップする」ということなのです。

![](https://storage.googleapis.com/zenn-user-upload/drawttp5tipfm1qkwadcsiyx6j5w)

## Readメソッド

`net.TCPConn`型の`Read()`の中身を掘り下げます。

先述した通り、`net.TCPConn`型の実体は非公開構造体`conn`です。そのため、`conn`型の`Read`メソッドがそのまま`net.TCPConn`型の`Read`メソッドとして機能します。


`(c *TCPConn) Read`が定義されていなくても、内部フィールド構造体の`(c *conn) Read`がそのまま`TCPConn`型のメソッドとして機能する挙動のことをメソッド委譲といいます。


その`conn`型の`Read`メソッドは、内部ではフィールド`fd`(`netFD`型)の`Read`メソッドを呼んでいます。


``` go
func (c *conn) Read(b []byte) (int, error) {
    // (略)
    n, err := c.fd.Read(b)
    // (略)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.2/src/net/net.go#179>

`netFD`型の`Read()`メソッドの中身では、`pfd`フィールド(`poll.FD`型)の`Read`メソッドを呼んでいます。


``` go
func (fd *netFD) Read(p []byte) (n int, err error) {
    n, err = fd.pfd.Read(p)
    // (略)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.2/src/net/fd_posix.go#54>

この`poll.FD`型の`Read`メソッドというのは、前章のファイルI/Oでも出てきたものです。ここから先は通常ファイルのI/Oと同じく、対応したfdのファイルの中身を読み込むためのシステムコール`syscall.Read`につながります。\
\"everything-is-a-file\"思想の名の通り、ネットワークコネクションからのデータ読み取りも、OSの世界においてはファイルの読み取りと変わらず`read`システムコールで処理されるのです。

`net.TCPConn`型の`Read`メソッドの処理手順をまとめます。

1.  `net.conn`型の`Read`メソッドを呼ぶ
2.  1の中で`net.netFD`型の`Read`メソッドを呼ぶ
3.  2の中で`poll.FD`型の`Read`メソッドを呼ぶ
4.  3の中で`syscall.Read`メソッドを呼ぶ
5.  OSカーネルのシステムコールで読み込み処理

## Writeメソッド

`net.TCPConn`型の`Write()`メソッドのほうも`Read`メソッドと同様の流れで実装されています。

1.  `net.conn`型の`Write`メソッドを呼ぶ
2.  1の中で`net.netFD`型の`Write`メソッドを呼ぶ
3.  2の中で`poll.FD`型の`Write`メソッドを呼ぶ
4.  3の中で`syscall.Write`メソッドを呼ぶ
5.  OSカーネルのシステムコールで書き込み処理

![](https://storage.googleapis.com/zenn-user-upload/pt2qg55vm9759qmjuif9mc88mm0u)

# まとめ

前章・本章とファイル・ネットワークのI/Oについて取り上げました。\
しかし、I/Oする対象こそ違えど、内部的な構造は両方とも

-   fdがある(=ファイルへのI/Oと見れる)
-   `Read()`メソッド、`Write()`メソッドのシグネチャが同じ
-   裏でシステムコールread()/write()を呼んでいる

等々、似ているところがあります。

次章では、これらI/Oをまとめてひっくるめて扱う抽象化の手段を紹介します。




# ioパッケージによる抽象化

# はじめに

今まで紹介してきたI/O読み書きメソッドは全て以下の形でした。


``` go
// バイトスライスpを用意して、そこに読み込んだ内容をいれる
Read(p []byte) (n int, err error)

// バイトスライスpの中身を書き込む
Write(p []byte) (n int, err error)
```


そのため、「ファイルでもネットワークでも何でもいいから、とにかく読み書きできるもの」が欲しい！というときに備えて、Goでは`io`パッケージによってインターフェース群が提供されています。

本章では

-   `io.Reader`と`io.Writer`
-   `io`で読み書きが抽象化されると何が嬉しいのか\
    について解説します。

# io.Readerの定義

`io.Reader`が一体なんなんかというと、次のメソッドをもつ**インターフェース**として定義されています。


``` go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```


出典:pkg.go.dev - io#Reader

つまり、`io.Reader`というのは、「何かを読み込む機能を持つものをまとめて扱うために抽象化されたもの」なのです。\
これまで扱った`os.File`型と`net.Conn`型はこの`io.Reader`インターフェースを満たします。

# io.Writerの定義

`io.Writer`は、以下の`Write`メソッドをもつ**インターフェース**として定義されています。


``` go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```


出典:pkg.go.dev - io#Writer

`io.Writer`は`io.Reader`と同様に、「何かに書き込む機能を持つものをまとめて扱うために抽象化されたもの」です。\
`os.File`型と`net.Conn`型はこの`io.Writer`インターフェースを満たします。

# 抽象化すると嬉しい具体例

「読み込み・書き込みを抽象化するようなインターフェースを作ったところで何が嬉しいの？」という方もいるでしょう。\
ここでは、`io`のインターフェースを利用して便利になる例を一つ作ってみます。

例えば、「どこかからの入力文字列を受け取って、その中の`Hello`を`Guten Tag`に置換する」という操作の実装を考えます。\
これを`io.Reader`を使わずに実装するとなると、「入力がファイルからの場合」と「入力がネットワークからの場合」という風に、具体型に沿って実装をいくつも用意しなくてはなりません。


``` go
// ファイルの中身を読み込んで文字列置換する関数
func FileTranslateIntoGerman(f *os.File) {
    data := make([]byte, 300)
    len, _ := f.Read(data)
    str := string(data[:len])

    result := strings.ReplaceAll(str, "Hello", "Guten Tag")
    fmt.Println(result)
}

// ネットワークコネクションからデータを受信して文字列置換する関数
func NetTranslateIntoGerman(conn net.Conn) {
    data := make([]byte, 300)
    len, _ := conn.Read(data)
    str := string(data[:len])

    result := strings.ReplaceAll(str, "Hello", "Guten Tag")
    fmt.Println(result)
}
```


2つの関数の実装は、引数の型が違うだけでほとんど同じです。

ここで、`io.Reader`インターフェースを使用することによって、2つの関数を1つにまとめることができます。


``` go
func TranslateIntoGerman(r io.Reader) {
    data := make([]byte, 300)
    len, _ := r.Read(data)
    str := string(data[:len])

    result := strings.ReplaceAll(str, "Hello", "Guten Tag")
    fmt.Println(result)
}
```


`io.Reader`インターフェース型の変数には、`os.File`型も`net.Conn`型も代入可能です。\
そのため、この`TranslateIntoGerman()`関数は、入力がファイルでもコネクションでも、どちらでも対応できる汎用性のある関数になりました。\
これがインターフェースによる抽象化のメリットです。

# まとめ

ここまで「`io`パッケージのインターフェースたちによって、どこからのI/Oであっても同様に扱える」ということをお見せしました。

次章からは、この`io.Reader`,
`io.Writer`を使った/に絡んだ便利なパッケージを紹介していきます。




# bufioパッケージによるbuffered I/O

# はじめに

標準パッケージの中にbufioパッケージというものがあります。\
ドキュメントによると、「bufioパッケージはbuffered
I/Oをやるためのもの」\[1\]と書かれています。

> Package bufio implements buffered I/O.\
> 出典:pkg.go.dev - bufio

これは普通にI/Oと一体何が違うのでしょうか。\
使い方と一緒に解説していきます。

# buffered I/O

## bufio.Reader型の特徴

`bufio`パッケージにはこのパッケージ特有の`bufio.Reader`型が存在します。\
`NewReader`関数を用いることで、`io.Reader`型から`bufio.Reader`型を作ることができます。


``` go
func NewReader(rd io.Reader) *Reader
```


出典:pkg.go.dev - bufio#NewReader


`bufio.Reader`型を作るための元になるリーダーが、具体型ではなく`io.Reader`であることで、「ネットワークでもファイルでもその他のI/Oであっても、buffered
I/Oにできる」のです。\
ここからも「`io`パッケージによるI/O抽象化」のメリットを感じることができます。


作った`bufio.Reader`は、普通の`io.Reader`とは何が違うのでしょうか。中身を見てみましょう。


``` go
type Reader struct {
    buf          []byte
    rd           io.Reader // reader provided by the client
    r, w         int       // buf read and write positions
    err          error
    lastByte     int // last byte read for UnreadByte; -1 means invalid
    lastRuneSize int // size of last rune read for UnreadRune; -1 means invalid
}
```


出典:<https://go.googlesource.com/go/+/go1.16.2/src/bufio/bufio.go#32>

ここで重要なのは`NewReader`関数の引数として与えられた`io.Reader`を格納する`rd`フィールドがあるということではなく、バイト列の`buf`フィールドがあるということです。


このバイト列`buf`の長さは、デフォルトでは`defaultBufSize = 4096`という定数で指定されています。


この`buf`がどんな役割を果たしているのでしょうか。それは`Read(p []byte)`メソッドの実装を見ればわかります。

1.  `len(p)`が内部バッファのサイズより大きい場合、読み込み結果を直接`p`にいれる
2.  `len(p)`が内部バッファのサイズより小さい場合、読み込み結果を一回内部バッファ`buf`に入れてから、その中身を`p`にコピー

出典:<https://go.googlesource.com/go/+/go1.16.2/src/bufio/bufio.go#198>

このように、ある特定条件下においては、「読み込んだ中身を内部バッファ`buf`に貯める」という動作が行われます。\
そのため、「もう変数`p`に内容を書き込み済みのデータも、`bufio.Reader`の内部バッファには残っている」状態になります。

## bufio.Writer型の特徴

`bufio.Reader`があるなら`bufio.Writer`も存在します。\
作り方も`bufio.Reader`と同様に、`io.Writer`型を`NewWriter`関数に渡すことで作ります。


``` go
func NewWriter(w io.Writer) *Writer
```


出典:pkg.go.dev - bufio#NewWriter

こうして作った`bufio.Writer`にも、内部バッファ`buf`が存在します。


``` go
type Writer struct {
    err error
    buf []byte
    n   int
    wr  io.Writer
}
```


出典:<https://go.googlesource.com/go/+/go1.16.2/src/bufio/bufio.go#558>

`Write(p []byte)`メソッドが実装されるときに、この内部バッファ`buf`がどう動くのでしょうか。\
実際に実装を確認してみると、以下のようになっています。

\<`p`の中身が全て処理されるまでこれを繰り返す>

1.  `len(p)`が内部バッファの空きより小さい場合(=`p`の中身を`buf`に書き込んでも`buf`に空きが余る場合)
    -   `p`の中身を一旦`buf`に書き込んでおく
2.  `len(p)`が内部バッファの空きより大きい場合(=`p`の中身を一旦全部`buf`に書き込むだけの余裕がない場合)
    -   `buf`が先頭から空いているなら、`p`の中身を直接メモリに書き込む(=`buf`を使わない)
    -   `buf`の空きが先頭からじゃないなら、
        1.  `buf`に入るだけデータを埋める
        2.  `buf`の中身をメモリに書き込む\[2\]
        3.  `p`の中で`buf`に書き込み済みのところを切る

つまり、「実際にデータをメモリに書き込むのは、内部バッファ`buf`の中身がいっぱいになったときのみ」という挙動をします。

わざわざこんな面倒なことをする理由に、OSがメモリを管理する方法が関連しています。\
基本的にOSは、**ブロック**単位(4KBだったり8KBだったりものにより様々)でメモリを割り当てています。\
そのため、「1byteの書き込みを4096回」と「4096byte(=4KB)の書き込みを1回」だったら後者の方が早いのです。

ユーザースペースでバッファリングすることで、中途半端な長さの書き込みを全て「ブロック単位の長さの書き込み」に整形することができるので、処理速度をあげることができるのです。

## ベンチマークによる実行時間比較

本当に`bufio`パッケージを使うことでI/Oが早くなっているのでしょうか。

> Measure. Don\'t tune for speed until you\'ve measured, and even then
> don\'t unless one part of the code overwhelms the rest.\
> (よく言われる意訳) 推測するな、計測しろ\
> 出典:Rob Pike\'s 5 Rules of
> Programming

こんな言葉もあることですし、それに従い実際にベンチマークをとって検証してみましょう。

検証環境は以下のものを使用しました。


    goos: darwin
    goarch: amd64
    cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz


### Readメソッド

まずは`io.Reader`と`bufio.Reader`型の`Read`メソッドを検証します。

以下のような関数を用意しました。


``` go
// サイズがFsizeのファイルをnbyteごと読む関数
func ReadOS(r io.Reader, n int) {
    data := make([]byte, n)

    t := Fsize / n
    for i := 0; i < t; i++ {
        r.Read(data)
    }
}
```


そして、ベンチマーク用のテストコードを以下のように書きました。


``` go
// ただのio用
func BenchmarkReadX(b *testing.B) {
    f, _ := os.Open("read.txt")
    defer f.Close()

    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        ReadOS(f, X)
    }
}

// bufio用
func BenchmarkBReadX(b *testing.B) {
    f, _ := os.Open("read.txt")
    defer f.Close()
    bf := bufio.NewReader(f)

    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        ReadOS(bf, X)
    }
}
```


ベンチマーク関数の名前`BenchmarkYReadX()`の名前は

-   `Y`: なしなら普通の`io`, `B`なら`bufio`での検証
-   `X`: `X`byteごとにファイルを読み込んでいく

です。

`go test -bench .`でのテスト結果は、以下のようになりました。


    BenchmarkRead1-8               1        1575492668 ns/op
    BenchmarkRead32-8             21          51526989 ns/op
    BenchmarkRead256-8           181           5954220 ns/op
    BenchmarkRead4096-8         3544            338707 ns/op

    BenchmarkBRead1-8             79          14302113 ns/op
    BenchmarkBRead32-8          1071          39197576 ns/op
    BenchmarkBRead256-8         1306           5104346 ns/op
    BenchmarkBRead4096-8        3427            373660 ns/op


1byteごと書き込んでいる場合、bufioの有無で110倍ものパフォーマンス差が生まれる結果となりました。

### Writeメソッド

今度は`io.Reader`と`bufio.Reader`型の`Write`メソッドを検証します。

検証用として以下のような関数を用意しました。


``` go
// サイズBsize分のデータを、nbyteごとに区切って書き込む
func WriteOS(w io.Writer, n int) {
    data := []byte(strings.Repeat("a", n))

    t := Bsize / n
    for i := 0; i < t; i++ {
        w.Write(data)
    }
}
```


そして、ベンチマーク用のテストコードを以下のように書きました。


``` go
// ただのio用
func BenchmarkWriteX(b *testing.B) {
    f, _ := os.Create("write.txt")
    defer f.Close()

    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        WriteOS(f, X)
    }
}

// bufio用
func BenchmarkBWriteX(b *testing.B) {
    f, _ := os.Create("write6.txt")
    defer f.Close()
    bf := bufio.NewWriter(f)

    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        WriteOS(bf, X)
    }
    bf.Flush()
}
```


`go test -bench .`でのテスト結果は、以下のようになりました。


    BenchmarkWrite1-8                    117          10157577 ns/op
    BenchmarkWrite32-8                  3280            330840 ns/op
    BenchmarkWrite256-8                27649             49118 ns/op
    BenchmarkWrite4096-8              206610              6637 ns/op

    BenchmarkBWrite1-8                 39537             29841 ns/op
    BenchmarkBWrite32-8               232269              5700 ns/op
    BenchmarkBWrite256-8              255998              5996 ns/op
    BenchmarkBWrite4096-8             193617              7128 ns/op


1byteごと読み込んでいる処理の場合、bufio使用なし/ありでそれぞれ10157577ns/29841nsと、約340倍ものパフォーマンスの差が出る結果となりました。\
読み込み単位のバイト数を増やすごとにパフォーマンス差はなくなっていきますが、それを抜きにしてもユーザースペースでのバッファリングの威力がよくわかる結果です。

# bufio.Scanner

`bufio`パッケージには、`Reader`とは別に`bufio.Scanner`という読み込みのための構造体がもう一つ存在します。\
`bufio.Reader()`での読み込みが「指定した長さのバイト列ごと」なのに対して、これは「トークンごとの読み込み」をできるようにすることで利便性を向上させたものです。

この章では`bufio.Scanner`について詳しくみていきます。

## トークン

### トークンとその利便性

`bufio.Scanner`で可能になる「トークン」ごとの読み取りですが、これは例えば

-   単語ごと(=スペース区切り)に読み取りたい
-   行ごと(=改行文字区切り)に読み取りたい

といった状況のときに威力を発揮します。\
上2つの例の場合、それぞれ「単語(word)」と「行(line)」をトークンにした`bufio.Scanner`を用意することで簡単に実現可能です。

これを`bufio.Reader`でやろうとすると、トークンごとの長さが時と場合によって変わるので、「まずは1000byte読み込んで、そこから単語や行ごとに区切って......」といった複雑な処理を自前で書かなくてはいけなくなります。\
`bufio.Scanner`はこの面倒な処理からユーザーを開放してくれます。

### トークン定義

トークンの定義は、`bufio`パッケージ内の`SplitFunc`型で行います。


``` go
type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
```


> SplitFunc is the signature of the split function used to tokenize the
> input.\
> (訳)`SplitFunc`型は、入力をトークンに分割するために使用する関数シグネチャです。\
> 出典:pkg.go.dev -
> bufio#SplitFunc

この`SplitFunc`型に代入することができる関数が、`bufio`内では4つ定義されています。


``` go
func ScanBytes(data []byte, atEOF bool) (advance int, token []byte, err error)
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error)
func ScanRunes(data []byte, atEOF bool) (advance int, token []byte, err error)
func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error)
```


つまり、`bufio`でデフォルトで定義されているトークンは以下の4つです。

-   バイトごと
-   行ごと
-   ルーンごと
-   単語ごと


「型リテラル`func ([]byte, bool) (int, []byte, error)`型の変数を`SplitFunc`型に代入できるの？違う型なのに？」と思った方は鋭いです。

実はこれは可能です。Goの言語仕様書で定義されている「代入可能性」には、「代入する変数と値の型が同一であること」という要項があります。\
今回の場合、`SplitFunc`というdefined
typeと型リテラル`func ([]byte, bool) (int, []byte, error)`は、underlying
typeが一緒なので型が同一判定されます。

Go
Playgroundで挙動を試してみた結果がこちらです。


## Scanner構造体について

### 内部構造

`bufio.Scanner`の内部構造は以下のようになっています。


``` go
type Scanner struct {
    r            io.Reader // The reader provided by the client.
    split        SplitFunc // The function to split the tokens.
    token        []byte    // Last token returned by split.
    buf          []byte    // Buffer used as argument to split.
    // (以下略)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/bufio/scan.go#30>

`bufio.Reader`型と同様に、内部にバッファを持っていることがわかります。\
つまり、`bufio.Scanner`の利用の裏ではbuffered I/Oが行われているのです。

また、`split`フィールドには、トークンを定義する`SplitFunc`型関数が格納されており、これに従ってスキャナーはトークン分割処理を行います。

scanner内では、tokenごとに区切る`SplitFunc`型の関数を内部に持っている。\
それをセットするのが`split()`メソッド。デフォルトはlineで区切られるようになっている。

### スキャナーの作成

`bufio.Scanner`の作成は、`bufio.Reader`の作成と同様に、`io.Reader`を引数に渡す`NewScanner`関数で行います。


``` go
func NewScanner(r io.Reader) *Scanner
```


出典:pkg.go.dev -
bufio#NewScanner

これで作成されたスキャナーは、デフォルトで「行」をトークンにするように設定されています。\
変更したい場合は、`Split`メソッドを使います。


``` go
// 引数で渡したSplitFuncでトークンを作る
func (s *Scanner) Split(split SplitFunc)
```


出典:pkg.go.dev -
bufio#Scanner.Split

## Scannerを使ってデータを読み取る

スキャナーを使ってデータを読みとるためには、「`Scan()`メソッドで読み込み→`Text()`メソッドで結果を取り出す」という手順を踏みます。

例えば、以下のようなテキストファイルを用意します。


    apple
    bird flies.
    cat is sleeping.
    dog


これを行ごとに読み取る処理を実装するには、例えば以下のようになります。


``` go
func main() {
    // ファイル(io.Reader)を用意
    f, _ := os.Open("text.txt")
    defer f.Close()

    // スキャナを用意(トークンはデフォルトの行のまま)
    sc := bufio.NewScanner(f)

    // EOFにあたるまでスキャンを繰り返す
    for sc.Scan() {
        line := sc.Text() // スキャンした内容を文字列で取得
        fmt.Println(line)
    }
}

/*
出力結果

apple
bird flies.
cat is sleeping.
dog
*/
```


## (余談)Scannerを使わなきゃ困っちゃう場面

実はこの`bufio.Scanner`、Goで競プロをやっている方なら馴染みがある概念ではないでしょうか。

### fmt.ScanでTLEが出ちゃう問題

競技プログラミングの問題において、大量のデータの入力が必要になる場合が存在します。\
例えば、このAtCoder Beginner Contest 144 -
E問題は、`2*N+2`個ものの数字が以下のように与えられます。


    N K
    A1 A2 ... AN
    F1 F2 ... FN


この問題の場合、`N`の最大値は`2*10^5`なので結構な数の入力があることになります。\
そのため、`fmt.Scan`を使うとTLE(時間切れによる不正解)判定が出てしまいます。


``` go
// TLEになったコードの断片
var N, K int
fmt.Scan(&N, &K)
 
A := make([]int, N)
for i := 0; i < N; i++ {
    fmt.Scan(&A[i])
}
F := make([]int, N)
for i := 0; i < N; i++ {
    fmt.Scan(&F[i])
}
```


### 解決策

これの解決策が`bufio.Scanner`の使用です。\
以下のようなコードはGoで競プロやるかたにとってはテンプレなのではないでしょうか。


``` go
var sc = bufio.NewScanner(os.Stdin)

func scanInt() int {
    sc.Scan()
    i, err := strconv.Atoi(sc.Text())
    if err != nil {
        panic(err)
    }
    return i
}

func main() {
    sc.Split(bufio.ScanWords)
    // (以下略)
}
```


これをファイル冒頭に入れておくことで、この問題での入力処理は以下のように書き換えられます。


``` go
N, K := scanInt(), scanInt()
 
A, F := make([]int, N), make([]int, N)
for i := 0; i < N; i++ {
    A[i] = scanInt()
}

for i := 0; i < N; i++ {
    F[i] = scanInt()
}
```


筆者はこの修正を施すことで無事にAC(正解)を通しました。

# まとめ

以上、`bufio`パッケージによるbuffered I/Oについて掘り下げました。

次章では、最後の競プロでも出てきた`fmt`での標準入力・出力について掘り下げます。



脚注


1.  
    \"bufio\"の名前の由来はおそらく\"buffer\"のbufと\"I/O\"のioを足したものでしょう
    

2.  
    この動作をflushといいます
    

3.  
    Go言語の父と呼ばれている人です
    





# fmtで学ぶ標準入力・出力

# はじめに

普段何気なく行う標準入力・出力もI/Oの一種です。\
ファイルやネットワークではなく、ターミナルからの入力・出力というのは裏で一体何が起こっているのでしょうか。\
本章では、`fmt`パッケージのコードと絡めてそれを探っていきます。

# 標準入力・標準出力の正体

いきなり答えを言ってしまうと、標準入力・標準出力自体は`os`パッケージで以下のように定義されています。


``` go
var (
    Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
    Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
)
```


出典:pkg.go.dev - os#Variables

出てくるワードを説明します。

-   `os.NewFile`関数:
    第二引数にとった名前のファイルを、第一引数にとったfd番号で`os.File`型にする関数
-   `syscall.Stdin`:
    `syscall`パッケージ内で`var Stdin  = 0`と定義された変数
-   `syscall.Stdout`:
    `syscall`パッケージ内で`var Stdout = 1`と定義された変数

つまり、

-   標準入力: ファイル`/dev/stdin`をfd0番で開いたもの
-   標準出力: ファイル`/dev/stdout`をfd1番で開いたもの

であり、ターミナルを経由した入力・出力も通常のファイルI/Oと同様に扱うことができるのです。


標準入力・出力に割り当てるfd番号を0と1にするのは一種の慣例です。\
また、標準エラー出力は慣例的にfd2番になります。


# fmt.Print系統

それでは「ターミナルに標準出力する」という処理がどのように実装されているのか、`fmt.Println`を一例にとってみていきましょう。


``` go
func Println(a ...interface{}) (n int, err error) {
    return Fprintln(os.Stdout, a...)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/fmt/print.go#273>

内部的には`fmt.Fprintln`を呼んでいることがわかります。\
その`fmt.Fprintln`は「第一引数にとった`io.Writer`に第二引数の値を書き込む」という関数です。


``` go
func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
    p := newPrinter()
    p.doPrintln(a)
    n, err = w.Write(p.buf)
    p.free()
    return
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/fmt/print.go#262>

実装的には「第一引数にとった`io.Writer`の`Write`メソッドを呼んでいる」だけです。

`os.Stdout`は`os.File`型の変数なので、当然`io.Writer`インターフェースは満たしています。\
そのため、そこへの出力は「ファイルへの出力」と全く同じ処理となります。

「標準出力はファイルなのだから、そこへの処理もファイルへの処理と同じ」という、直感的にわかりやすい結果です。

# fmt.Scan系統

出力をみた後は、今度は標準入力のほうをみてみましょう。

今回掘り下げるのは`fmt.Scan`関数です。内部的にはこれは`fmt.Fscan`を呼んでいるだけです。


``` go
func Scan(a ...interface{}) (n int, err error) {
    return Fscan(os.Stdin, a...)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/fmt/scan.go#63>

ここで出てきた`fmt.Fscan`関数は、第一引数の`io.Reader`から読み込んだデータを第二引数に入れる関数です。\
内部実装は以下のようになっています。


``` go
func Fscan(r io.Reader, a ...interface{}) (n int, err error) {
    s, old := newScanState(r, true, false)  // newScanState allocates a new ss struct or grab a cached one.
    n, err = s.doScan(a)
    s.free(old)
    return
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/fmt/scan.go#121>

ざっくりと解説すると

1.  `newScanState`から得た変数`s`は、第一引数で渡した`io.Reader`(ここでは`os.Stdin`ファイル)を内包した構造体
2.  1で得た構造体の`s.doScan`メソッドの内部で、第一引数`r`の`Read`メソッドを呼んでいる

「標準入力はファイルなのだから、そこへの処理もファイルへの処理と同じ」という、標準出力と同様の結果になります。

# まとめ

ここでは、「標準入力・出力はファイル`/dev/stdin`・`/dev/stdout`への入出力と同じ」ということを取り上げました。

次章では、普段何気なく扱っているものを`io.Reader`/`io.Writer`として扱うための便利パッケージを紹介します。




# bytesパッケージとstringsパッケージ

# はじめに

`io.Reader`と`io.Writer`を満たす型として、`bytes`パッケージの`bytes.Buffer`型が存在します。\
また、`strings`パッケージの`strings.Reader`型は`io.Reader`を満たします。

本章では、これらの型について解説します。

# bytesパッケージのbytes.Buffer型

まずは、`bytes.Buffer`型の構造体の中身を確認してみましょう。


``` go
type Buffer struct {
    buf      []byte
    // (略)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/bytes/buffer.go#20>

構造として特筆すべきなのは、中身にバイト列を持っているだけです。\
これだけだったら「そのまま`[]byte`を使えばいいじゃないか」と思うかもしれませんが、パッケージ特有の型を新しく定義することによって、メソッドを好きに付けられるようになります。

というわけで、`bytes.Buffer`には`Read`メソッド、`Write`メソッドがついています。\
これによって`io.Reader`と`io.Writer`を満たすようになっています。

## Writeメソッド

`Write`メソッドは、レシーバーのバッファの「中に」データを書き込むためのメソッドです。

使用例を以下に示します。


``` go
// bytes.Bufferを用意
// (bytes.Bufferは初期化の必要がありません)
var b bytes.Buffer
b.Write([]byte("Hello"))

// バッファの中身を確認してみる
fmt.Println(b.String())

// (出力)
// Hello
```


参考:pkg.go.dev -
bytes#Buffer-Example

## Readメソッド

`Read`メソッドは、レシーバーバッファから「中を」読み取るためのメソッドです。


``` go
// 中にデータを入れたバッファを用意
var b bytes.Buffer
b.Write([]byte("World"))

// plainの中にバッファの中身を読み取る
plain := make([]byte, 10)
b.Read(plain)

// 読み取り後のバッファの中身と読み取り結果を確認
fmt.Println("buffer: ", b.String())
fmt.Println("output:", string(plain))

// buffer:  
// output: World
```


バッファの中からは`World`というデータが見えなくなり、きちんと変数`plain`に読み込みが成功しています。

# stringsパッケージのstrings.Reader型

`strings`パッケージは、文字列を置換したり辞書順でどっちが先か比べたりという単なる便利屋さんだけではないのです。\
`bytes.Buffer`型と同じく、文字列型をパッケージ独自型でカバーすることで、`io.Reader`に代入できるようにした型も定義されているのです。

そんな独自型`strings.Reader`型は、構造体内部に文字列を内包しています。


``` go
type Reader struct {
    s        string
    // (略)
}
```


出典:<https://go.googlesource.com/go/+/go1.16.3/src/strings/reader.go#17>

これは`Read`メソッドをもつ、`io.Reader`インターフェースを満たす構造体です。


`strings.Reader`型に`Write`メソッドはないので、`io.Writer`は満たしません。


## Readメソッド

`Read`メソッドは、文字列から「中を」読み取るためのメソッドです。\
使用例を示します。


``` go
// NewReader関数から
// strings.Reader型のrdを作る
str := "Hellooooooooooooooooooooooooooo!"
rd := strings.NewReader(str)

// rowの中に読み取り
row := make([]byte, 10)
rd.Read(row)

// 読み取り結果確認
fmt.Println(string(row)) // Helloooooo
```


# これらの使いどころ

「バイト列や文字列を`io.Reader`・`io.Writer`に入れられるようにしたところで何が嬉しいの？」という疑問を持った方もいるかと思います。\
ここからはそんな疑問に対して、ここで紹介した型の使い所を一つ紹介したいと思います。

## テストをスマートに書く

`io`の章で書いた`TranslateIntoGerman`関数を思い出してください。


``` go
// 引数rで受け取った中身を読み込んで
// Hello → Guten Tagに置換する関数
func TranslateIntoGerman(r io.Reader) string {
    data := make([]byte, 300)
    len, _ := r.Read(data)
    str := string(data[:len])

    result := strings.ReplaceAll(str, "Hello", "Guten Tag")
    return result
}
```


この関数のテストを書くとき、皆さんならどうするでしょうか。\
「`io.Reader`を満たすものしか引数に渡せない......テストしたい内容が書いてあるファイルを1個1個用意するか...？」と思ったこともいるでしょう。

ですが「ファイルを1個1個用意する」とかいう面倒な方法をせずとも、`strings.Reader`型を使うことで、テスト内容をコード内で用意することができるのです。


``` go
func Test_Translate(t *testing.T) {
    // テストしたい内容を文字列ベースで用意
    tests := []struct {
        name string
        arg  string
        want string
    }{
        {
            name: "normal",
            arg:  "Hello, World!",
            want: "Guten Tag, World!",
        },
        {
            name: "repeat",
            arg:  "Hello World, Hello Golang!",
            want: "Guten Tag World, Guten Tag Golang!",
        },
    }

    // TranslateIntoGerman関数には
    // strings.NewReader(tt.args)で用意したstrings.Reader型を渡す
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := TranslateIntoGerman(strings.NewReader(tt.arg)); got != tt.want {
                t.Errorf("got %v, but want %v", got, tt.want)
            }
        })
    }
}
```





# おわりに

# おわりに

というわけで、GoでI/Oに関わるものを片っ端から書き連ねました。\
完全にごった煮状態の本ですがいかがでしたでしょうか。

I/Oは、根本を理解しようとすると低レイヤの知識まで必要になってくるのでなかなか難しいですが、この本が皆さんの理解の一助になっていれば幸いです。

コメントによる編集リクエスト・情報提供は大歓迎ですので、どしどし書き込んでいってください。\
連絡先: 作者Twitter \@saki_engineer

# 参考文献

## 書籍 Linux System Programming






<https://learning.oreilly.com/library/view/linux-system-programming/0596009585/>

オライリーの本です。\
Linuxでの低レイヤ・カーネル内部まわりの話がこれでもかというほど書かれています。\
今回この本を執筆するにあたって、1\~4章のI/Oの章を大いに参考にしました。

## 書籍 Software Design 2021年1月号






<https://gihyo.jp/magazine/SD/archive/2021/202101>

この本のGo特集第2章が、tenntennさん(\@tenntenn)が執筆された`io`章です。\
このZenn本では`io.Reader`と`io.Writer`しか取り上げませんでしたが、Software
Designの記事の方には他の`io`便利インターフェースについても言及があります。

## Web連載 Goならわかるシステムプログラミング






<https://ascii.jp/serialarticles/1235262/>

渋川よしきさん(\@shibu_jp)が書かれたWeb連載です。\
Goの視点からみた低レイヤの話がとても詳しく書かれています。

以下の回を大いに参考にしました。

-   第2回
    低レベルアクセスへの入り口（1）：io.Writer
-   第3回
    低レベルアクセスへの入り口（2）：io.Reader前編
-   第4回
    低レベルアクセスへの入り口（3）：io.Reader後編
-   第5回
    Goから見たシステムコール
-   第6回
    GoでたたくTCPソケット（前編）
-   第7回
    GoでたたくTCPソケット（後編）

## Qiita記事 Go言語を使ったTCPクライアントの作り方






<https://qiita.com/tutuz/items/e875d8ea3c31450195a7>

Go Advent Calender 2020 10日目にTsuji
Daishiro(\@d_tutuz)さんが書かれた記事です。\
TCPネットワークにおけるシステムコールは、この本ではsocket()しか取り上げませんでしたが、この記事ではさらに詳しいところまで掘り下げています。

## GopherCon 2019: Dave Cheney - Two Go Programs, Three Different Profiling Techniques

動画\







# エラーが発生しました。


www.youtube.com
での動画の視聴をお試しください。また、お使いのブラウザで JavaScript
が無効になっている場合は有効にしてください。





<https://www.youtube.com/watch?v=nok0aYiGiYA>\
サマリー記事\






<https://about.sourcegraph.com/go/gophercon-2019-two-go-programs-three-different-profiling-techniques-in-50-minutes/>

Dave
Cheneyさん(\@davecheney)によるGoCon2019のセッション(英語)です。\
前半部分が「ユーザースペースでバッファリングしたI/Oは早いぞ」という内容です。\
セッション中に実際にコードを書いて、それを分析ツールでどこが遅いのかを確かめながらコードを改善していく様子がよくわかります。

## 記事 How to read and write with Golang bufio






<https://www.educative.io/edpresso/how-to-read-and-write-with-golang-bufio>

`bufio.Writer`を使った際の内部バッファの挙動がイラスト付きでわかりやすく書かれています。




