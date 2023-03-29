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

TODO:

Instead of ast.Inspect and direct looping, find a way to get all the `//go:generate ` in the file and spit out one generated file.

```go

// something like
/*
   var files []*ast.File
   typeMap := make(map[string]*ast.TypeSpec)

   for _, pkg := range pkgs {
        for _, file := range pkg.Syntax {
                for _, decl := range file.Decls {
                        genDecl, ok := decl.(*ast.GenDecl)
                        if !ok || genDecl.Tok != token.TYPE {
                           continue
                        }        

           for _, spec := range genDecl.Specs {
                   if typeSpec, ok := spec.(*ast.TypeSpec); ok {
                           typeName := typeSpec.Name.String()
                           if contains(typeNameList, typeName) {
                                   typeMap[typeName] = typeSpec
                                   files = append(files, file)
                           }
                   }
           }

        }
   }
   }
*/
```
