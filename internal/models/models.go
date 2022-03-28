package models

import (
	"strings"
	"unicode"
)

type Expression struct {
	NonQualifiedValue string
	IsStar            bool
	IsVariadic        bool
	IsWriter          bool
	TestOnlyPackage   bool
	Underlying        string
	Package           string
}

func (e *Expression) IsBasicType() bool {
	return isBasicType(e.NonQualifiedValue) || isBasicType(e.Underlying)
}

func (e *Expression) IsWellKnownType() bool {
	return e.IsBasicType() || isWellKnownType(e.NonQualifiedValue) || isWellKnownType(e.Underlying)
}

func (e *Expression) HasPackagePath() bool {
	return strings.Contains(e.NonQualifiedValue, ".")
}

func (e *Expression) Value() string {
	if e.TestOnlyPackage && !e.IsWellKnownType() && !e.HasPackagePath() {
		return e.Package + "." + e.NonQualifiedValue
	}
	return e.NonQualifiedValue
}

func (e *Expression) String() string {
	value := e.Value()
	if e.IsStar {
		value = "*" + value
	}
	if e.IsVariadic {
		return "[]" + value
	}
	return value
}

type Field struct {
	Name  string
	Type  *Expression
	Index int
}

func (f *Field) IsWriter() bool {
	return f.Type.IsWriter
}

func (f *Field) IsStruct() bool {
	return strings.HasPrefix(f.Type.Underlying, "struct")
}

func (f *Field) IsBasicType() bool {
	return f.Type.IsBasicType()
}

func isBasicType(t string) bool {
	switch t {
	case "bool", "string", "int", "int8", "int16", "int32", "int64", "uint",
		"uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune",
		"float32", "float64", "complex64", "complex128":
		return true
	default:
		return false
	}
}

func isWellKnownType(t string) bool {
	switch t {
	case "error", "[]byte", "interface{}", "any":
		return true
	default:
		return false
	}
}

func (f *Field) IsNamed() bool {
	return f.Name != "" && f.Name != "_"
}

func (f *Field) ShortName() string {
	return strings.ToLower(string([]rune(f.Type.NonQualifiedValue)[0]))
}

type Receiver struct {
	*Field
	Fields []*Field
}

func (r *Receiver) Name() string {
	var n string
	if r.IsNamed() {
		n = r.Field.Name
	} else {
		n = r.ShortName()
	}
	if n == "name" {
		// Avoid conflict with test struct's "name" field.
		n = "n"
	} else if n == "t" {
		// Avoid conflict with test argument.
		// "tr" is short for t receiver.
		n = "tr"
	}
	return n
}

type Function struct {
	NonQualifiedName string
	IsExported       bool
	Receiver         *Receiver
	Parameters       []*Field
	Results          []*Field
	ReturnsError     bool
	TestOnlyPackage  bool
	Package          string
}

func (f *Function) TestParameters() []*Field {
	var ps []*Field
	for _, p := range f.Parameters {
		if p.IsWriter() {
			continue
		}
		ps = append(ps, p)
	}
	return ps
}

func (f *Function) TestResults() []*Field {
	var ps []*Field
	ps = append(ps, f.Results...)
	for _, p := range f.Parameters {
		if !p.IsWriter() {
			continue
		}
		ps = append(ps, &Field{
			Name: p.Name,
			Type: &Expression{
				NonQualifiedValue: "string",
				IsWriter:          true,
				Underlying:        "string",
			},
			Index: len(ps),
		})
	}
	return ps
}

func (f *Function) ReturnsMultiple() bool {
	return len(f.Results) > 1
}

func (f *Function) OnlyReturnsOneValue() bool {
	return len(f.Results) == 1 && !f.ReturnsError
}

func (f *Function) OnlyReturnsError() bool {
	return len(f.Results) == 0 && f.ReturnsError
}

func (f *Function) FullName() string {
	var r string
	if f.Receiver != nil {
		r = f.Receiver.Type.NonQualifiedValue
	}
	return strings.Title(r) + strings.Title(f.NonQualifiedName)
}

func (f *Function) TestName() string {
	if strings.HasPrefix(f.NonQualifiedName, "Test") {
		return f.NonQualifiedName
	}
	if f.Receiver != nil {
		receiverType := f.Receiver.Type.NonQualifiedValue
		if unicode.IsLower([]rune(receiverType)[0]) {
			receiverType = "_" + receiverType
		}
		return "Test" + receiverType + "_" + f.NonQualifiedName
	}
	if unicode.IsLower([]rune(f.NonQualifiedName)[0]) {
		return "Test_" + f.NonQualifiedName
	}
	return "Test" + f.NonQualifiedName
}

func (f *Function) IsNaked() bool {
	return f.Receiver == nil && len(f.Parameters) == 0 && len(f.Results) == 0
}

func (f *Function) Name() string {
	if f.Receiver != nil {
		return f.Receiver.Name() + "." + f.NonQualifiedName
	}
	if f.TestOnlyPackage && f.Package != "" {
		return f.Package + "." + f.NonQualifiedName
	}
	return f.NonQualifiedName
}

type Import struct {
	Name, Path string
}

type Package struct {
	Name     string
	TestOnly bool
}

func (p *Package) TestOnlyName() string {
	if strings.HasSuffix(p.Name, "_test") {
		return p.Name
	}
	return p.Name + "_test"
}

func (p *Package) String() string {
	if p.TestOnly {
		return p.TestOnlyName()
	}
	return p.Name
}

type Header struct {
	Comments []string
	Package  *Package
	Imports  []*Import
	Code     []byte
}

type Path string

func (p Path) TestPath() string {
	if !p.IsTestPath() {
		return strings.TrimSuffix(string(p), ".go") + "_test.go"
	}
	return string(p)
}

func (p Path) IsTestPath() bool {
	return strings.HasSuffix(string(p), "_test.go")
}
