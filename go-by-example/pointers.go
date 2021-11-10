package main

import "fmt"

// 値と対比しながら、ポインタがどのように動作するかを
// 2 つの関数 zeroval と zeroptr を使って示します。
// zeroval は int 型のパラメーターをもつので、引数は
// 値で渡されます。 zeroval は、呼び出し元の関数とは
// 別の ival コピーを受け取ります。
func zeroval(ival int) {
	ival = 0
}

// 一方、 zeroptr は *int 型のパラメーターをもつので、
// int のポインタを受け取ります。関数本体の *iptr は、
// ポインタを デリファレンス (dereferences) し、
// ポインタの指すメモリアドレスから現在の値を取得します。
// ポインタのデリファレンスに値を代入すると、 参照されて
// いるアドレスにある値が変わります。
func zeoptr(iptr *int) {
	*iptr = 0
}

func main() {
	i := 1
	fmt.Println("initial:", i)

	zeroval(i)
	fmt.Println("zeroval:", i)

	// &i という構文で、 i のメモリアドレス、
	// すなわち i へのポインタを取得できます。
	zeoptr(&i)
	fmt.Println("zeroptr:", i)

	// ポインタは表示することもできます。
	fmt.Println("pointer:", &i)
}
