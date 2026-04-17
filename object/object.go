package object

import (
	"fmt"
	"interpreter_go/ast"
	"bytes"
	"strings"
)

const (
	INTEGER_OBJ = "INTEGER"
	STRING_OBJ = "STRING"
	BOOLEAN_OBJ = "BOOLEAN"
	ERROR_OBJ = "ERROR"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ = "FUNCTION_OBJ"
	BUILTIN_OBJ = "BUILTIN"
	NULL_OBJ = "NULL"
)

type ObjectType string

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type String struct {
	Value string
}

type Boolean struct {
	Value bool
}

type Function struct {
	Parameters []*ast.Identifier
	Body *ast.BlockStatement
	Env *Environment
}


type ReturnValue struct {
	Value Object
}

type Error struct {
	Message string
}

type Null struct {}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func(b *Builtin) Inspect() string {
	return "builtin function"
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ 
}

func (f *Function) Inspect() string {
	var buff bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	buff.WriteString("fn")
	buff.WriteString("(")
	buff.WriteString(strings.Join(params, ", "))
	buff.WriteString(") {\n")
	buff.WriteString(f.Body.String())
	buff.WriteString("\n}")
	return buff.String()
}


