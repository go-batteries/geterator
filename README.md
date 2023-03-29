## geterator

This is something like a `stringer` in golang, which allows you to generate getters for a given type

Example usage:

```go
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

```

`$> go generate ./...`

This generates two files. `utils/foo_gen.go` and `utils/bar_gen.go`


Example generate config

```
// this is a generated file, please don't edit it by hand
package utils

func (c Bar) GetDrinks() []int {
	return c.Drinks
}
func (c Bar) GetNamer() Foo {
	return c.Namer
}
```

Right now the limitation is it creates one file per struct. Maybe I will find a way out.
