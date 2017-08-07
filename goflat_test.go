package gdl

import (
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"testing"

	"h12.me/gdl/internal/testpkg"
)

func TestGoflat(t *testing.T) {
	v := &testpkg.TestStruct{}
	pkg, err := Parse(v)
	if err != nil {
		t.Fatal(err)
	}
	pkg.ToFlatBuffers(os.Stdout)
}

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

type printer struct {
	w io.Writer
}

func (p printer) Printlnf(format string, v ...interface{}) {
	fmt.Fprintf(p.w, format, v...)
	fmt.Fprintln(p.w, "")
}

func (p *Package) ToFlatBuffers(w io.Writer) error {
	fp := printer{w}.Printlnf
	fp("namespace %s;", p.Name)
	fp("")
	for _, t := range p.Types {
		switch t.Kind {
		case Struct:
			fp("struct %s {", t.Name)
			for _, field := range t.Fields {
				fp("%s:%s;", field.Name, idlType(field.Type.Name))
			}
			fp("}")
			fp("")
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
	case reflect.Int, reflect.String:
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

func (p *Package) ToProtocolBuffers(w io.Writer) error {
	fp := printer{w}.Printlnf
	fp(`syntax = "proto3";`)
	fp(``)
	for _, t := range p.Types {
		switch t.Kind {
		case Struct:
			fp("struct %s {", t.Name)
			for _, field := range t.Fields {
				fp("%s:%s;", field.Name, idlType(field.Type.Name))
			}
			fp("}")
		}
	}
	return nil
}
