package utils

//go:generate go run github.com/go-batteries/geterator -type=Foo
type Foo struct {
	Name   string
	hidden bool
}

//go:generate go run github.com/go-batteries/geterator -type=Bar -private
type Bar struct {
	Drinks []int
	Namer  Foo
	hidden bool
}
