package object

import (
	"fmt"
	"interpreter_go/ast"
	"bytes"
	"strings"
	"hash/fnv"
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
	ARRAY_OBJ = "ARRAY"
	HASH_OBJ = "HASH"
)

type Hashable interface {
	HashKey() HashKey
}

type ObjectType string

type BuiltinFunction func(args ...Object) Object

type HashKey struct {
	Type ObjectType
	Value uint64
}

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

type Array struct {
	Elements []Object
}

type Error struct {
	Message string
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

type Null struct {}

type HashPair struct {
	Key Object
	Value Object
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	var buff bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	buff.WriteString("{")
	buff.WriteString(strings.Join(pairs, ", "))
	buff.WriteString("}")
	return buff.String()
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *Array) Inspect() string {
	var buff bytes.Buffer

	elements := []string{}

	for _, ele := range a.Elements {
		elements = append(elements, ele.Inspect())
	}
	buff.WriteString("[")
	buff.WriteString(strings.Join(elements, ","))
	buff.WriteString("]")

	return buff.String()
}

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

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

