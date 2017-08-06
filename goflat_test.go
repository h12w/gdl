package goflat

import (
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"testing"

	"h12.me/goflat/internal/testpkg"
)

type (
	Package struct {
		Name  string
		Types []*Type
	}
	Type struct {
		Name   string
		Kind   Kind
		Fields []Field
	}
	Field struct {
		Name string
		Type *Type
	}
	Kind int
)

const (
	Int    = Kind(reflect.Int)
	Struct = Kind(reflect.Struct)
)

func (p *Package) IDL(w io.Writer) error {
	fp := fmt.Fprintf
	fp(w, "namespace %s;\n\n", p.Name)
	for _, t := range p.Types {
		switch t.Kind {
		case Struct:
			fp(w, "struct %s {\n", t.Name)
			for _, field := range t.Fields {
				fp(w, "%s:%s;\n", field.Name, idlType(field.Type.Name))
			}
			fp(w, "}\n")
		}
	}
	return nil
}
func idlType(goType string) string {
	switch goType {
	case "int":
		return "long"
	}
	return goType
}

func TestGoflat(t *testing.T) {
	v := &testpkg.TestStruct{}
	pkg, err := Parse(v)
	if err != nil {
		t.Fatal(err)
	}
	pkg.IDL(os.Stdout)
}

func Parse(v interface{}) (*Package, error) {
	t := reflect.ValueOf(v).Type()
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pkg := &Package{
		Name: path.Base(t.PkgPath()),
	}
	_, err := pkg.parseType(t)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func (p *Package) parseType(t reflect.Type) (*Type, error) {
	for _, typ := range p.Types {
		if typ.Name == t.Name() {
			return typ, nil
		}
	}
	typ := &Type{
		Name: t.Name(),
		Kind: Kind(t.Kind()),
	}
	switch t.Kind() {
	case reflect.Int:
		return typ, nil
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fieldType, err := p.parseType(f.Type)
			if err != nil {
				return nil, err
			}
			typ.Fields = append(typ.Fields, Field{
				Name: f.Name,
				Type: fieldType,
			})
		}
	default:
		return nil, fmt.Errorf("unsupported kind %v", t.Kind())
	}
	p.Types = append(p.Types, typ)
	return typ, nil
}
