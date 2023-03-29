package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	typeName = flag.String("type", "", "struct names. must be present in go file. -type=Config")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("getter: ")

	flag.Parse()

	if len(*typeName) == 0 {
		log.Fatal("type name is missing")
	}

	cfg := &packages.Config{
		Mode:  packages.LoadSyntax,
		Tests: false,
	}

	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	var dir string

	if len(dirs) == 1 && isDirectory(dirs[0]) {
		dir = dirs[0]
	} else {
		dir = filepath.Dir(dirs[0])
	}

	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		log.Fatal(err)
	}

	var pkg *packages.Package

	for _, p := range pkgs {
		if p.Types == nil {
			continue
		}

		for _, f := range p.Syntax {
			ast.Inspect(f, func(node ast.Node) bool {
				// Look for a type definition with the given name.
				if typeSpec, ok := node.(*ast.TypeSpec); ok {
					if typeSpec.Name.Name == *typeName {
						pkg = p
						return false
					}
				}
				return true
			})
			if pkg != nil {
				break
			}
		}
		if pkg != nil {
			break
		}
	}

	if pkg == nil {
		log.Fatalf("error: type %q not found\n", *typeName)
	}

	var buf = bytes.Buffer{}
	buf.WriteString("// this is a generated file, please don't edit it by hand \n")
	buf.WriteString(fmt.Sprintf("package %s\n", pkg.Name))


	for _, f := range pkg.Syntax {
		ast.Inspect(f, func(node ast.Node) bool {
			// Look for a type definition with the given name.
			if typeSpec, ok := node.(*ast.TypeSpec); ok {
				if typeSpec.Name.Name == *typeName {
					// Generate getters for the fields of the struct.
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						for _, field := range structType.Fields.List {
							// Skip unexported fields.
							if field.Names == nil {
								continue
							}

							fieldName := field.Names[0].Name
							fieldType := exprToString(field.Type)

							buf.WriteString(
								fmt.Sprintf("func (c %s) Get%s() %s {\n", *typeName, strings.Title(fieldName), fieldType),
							)
							buf.WriteString(fmt.Sprintf("return c.%s\n", fieldName))
							buf.WriteString("}\n")
						}
					}

					return false
				}
			}

			return true
		})
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	filename := fmt.Sprintf("%s_gen.go", strings.ToLower(*typeName))
	path := filepath.Join(dir, filename)
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	defer file.Close()

	if _, err := file.Write(out); err != nil {
		log.Fatal(err)
	}
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	case *ast.ArrayType:
		return "[]" + exprToString(e.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", exprToString(e.Key), exprToString(e.Value))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", e.X, e.Sel)
	case *ast.StructType:
		return "struct{}"
	default:
		log.Fatal("unsupported type")
	}

	return "interface{}"
}

// containsGenerateDirective returns true if the specified file contains a go:generate directive.
func containsGenerateDirective(pkg *packages.Package, filename string) bool {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("error parsing file: %v", err)
	}
	for _, commentGroup := range node.Comments {
		for _, comment := range commentGroup.List {
			if strings.HasPrefix(comment.Text, "//go:generate") {
				return true
			}
		}
	}
	return false
}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}

	return info.IsDir()
}
