package utils

//go:generate go run github.com/go-batteries/geterator -type=Foo
type Foo struct {
	Name string
}

//go:generate go run github.com/go-batteries/geterator -type=Bar
type Bar struct {
	Drinks []int
	Namer  Foo
}
