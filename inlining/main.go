// go build -gcflags=-m

package main

type foo struct {
	a string
	b string
	c string
}

func newSugar(a, b, c string) *foo {
	return &foo{
		a: a,
		b: b,
		c: c,
	}
}

func newUgly(a, b, c string) *foo {
	return &foo{a: a, b: b, c: c}
}

func main() {

}
