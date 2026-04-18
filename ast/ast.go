package ast

import (
	"bytes"
	"interpreter_go/token"
	"strings"
	"fmt"
)


type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type ArrayLiteral struct {
	Token token.Token
	Elements []Expression
}

type LetStatement struct {
	Token token.Token
	Name *Identifier 
	Value Expression
}

type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

type StringLiteral struct {
	Token token.Token
	Value string
}

type InfixExpression struct {
	Token token.Token
	Right Expression
	Operator string
	Left Expression
}

type PrefixExpression struct {
	Token token.Token
	Operator string
	Right Expression
}

type Boolean struct {
	Token token.Token
	Value bool
}

type Identifier struct {
	Token token.Token
	Value string
}

type IfExpression struct {
	Token token.Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

type BlockStatement struct {
	Token token.Token
	Statements []Statement
}

type FunctionLiteral struct {
	Token token.Token
	Parameters []*Identifier
	Body *BlockStatement
}

type CallExpression struct {
	Token token.Token
	Function Expression
	Arguments []Expression
}

type IndexExpression struct {
	Token token.Token
	Left Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal}
func (ie *IndexExpression) String() string {
	var buff bytes.Buffer

	buff.WriteString("(")
	buff.WriteString(ie.Left.String())
	buff.WriteString("[")
	buff.WriteString(ie.Index.String())
	buff.WriteString("])")

	return buff.String()
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {

	var buff bytes.Buffer

	eles := []string{}

	for _, ele := range al.Elements {
		eles = append(eles, ele.String())
	}

	buff.WriteString("[")
	buff.WriteString(strings.Join(eles, ","))
	buff.WriteString("]")

	return buff.String()
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var buff bytes.Buffer

	buff.WriteString(ls.TokenLiteral() + " ")
	buff.WriteString(ls.Name.String())
	buff.WriteString(" = ")

	if ls.Value != nil {
		buff.WriteString(ls.Value.String())
	}

	buff.WriteString(";")

	return buff.String()
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value }

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { return il.Token.Literal }

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string { return sl.Token.Literal }

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var buff bytes.Buffer
	
	buff.WriteString("(")
	buff.WriteString(pe.Operator)
	buff.WriteString(pe.Right.String())
	buff.WriteString(")")

	return buff.String()
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var buff bytes.Buffer
	
	buff.WriteString("(")
	buff.WriteString(ie.Left.String())
	buff.WriteString(" " + ie.Operator + " ")
	buff.WriteString(ie.Right.String())
	buff.WriteString(")")

	return buff.String()
}

func (ife *IfExpression) expressionNode() {}
func (ife *IfExpression) TokenLiteral() string { return ife.Token.Literal }
func (ife *IfExpression) String() string {
	var buff bytes.Buffer
	
	buff.WriteString("if ")
	buff.WriteString(ife.Condition.String())
	buff.WriteString(" ")
	buff.WriteString(ife.Consequence.String())
	
	if ife.Alternative != nil {
		buff.WriteString(" else ")
		buff.WriteString(ife.Alternative.String())
	}

	return buff.String()
}


func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var buff bytes.Buffer
	params := []string{}
	
	for _, p := range fl.Parameters {
	params = append(params, p.String())
	}
	buff.WriteString(fl.TokenLiteral())
	buff.WriteString("(")
	buff.WriteString(strings.Join(params, ", "))
	buff.WriteString(") ")
	buff.WriteString(fl.Body.String())
	
	return buff.String()
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var buff bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	buff.WriteString(ce.Function.String())
	buff.WriteString("(")
	buff.WriteString(strings.Join(args, ", "))
	buff.WriteString(")")
	return buff.String()
}

func (bs *BlockStatement) expressionNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var buff bytes.Buffer
	
	for _, s := range bs.Statements {
		buff.WriteString(s.String())
	}

	return buff.String()
}


func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string { return b.Token.Literal }

func (p *Program) String() string {
	var buff bytes.Buffer
	for _, s := range p.Statements {
		buff.WriteString(s.String())
	}

	return buff.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func PrintAST(node Node, indent string) {
	switch n := node.(type) {
	case *Program:
		fmt.Printf("%sProgram:\n", indent)
		for _, stmt := range n.Statements {
			PrintAST(stmt, indent+"  ")
		}
	case *LetStatement:
		fmt.Printf("%sLetStatement:\n", indent)
		fmt.Printf("%s  Name: %s\n", indent, n.Name.Value)
		fmt.Printf("%s  Value:\n", indent)
		if n.Value != nil {
			PrintAST(n.Value, indent+"    ")
		}
	case *ReturnStatement:
		fmt.Printf("%sReturnStatement:\n", indent)
		if n.ReturnValue != nil {
			fmt.Printf("%s  ReturnValue:\n", indent)
			PrintAST(n.ReturnValue, indent+"    ")
		}
	case *ExpressionStatement:
		fmt.Printf("%sExpressionStatement:\n", indent)
		if n.Expression != nil {
			PrintAST(n.Expression, indent+"  ")
		}
	case *InfixExpression:
		fmt.Printf("%sInfixExpression (%s):\n", indent, n.Operator)
		fmt.Printf("%s  Left:\n", indent)
		PrintAST(n.Left, indent+"    ")
		fmt.Printf("%s  Right:\n", indent)
		PrintAST(n.Right, indent+"    ")
	case *IntegerLiteral:
		fmt.Printf("%sIntegerLiteral: %d\n", indent, n.Value)
	case *Identifier:
		fmt.Printf("%sIdentifier: %s\n", indent, n.Value)
	case *Boolean:
		fmt.Printf("%sBoolean: %t\n", indent, n.Value)
	case *IfExpression:
		fmt.Printf("%sIfExpression:\n", indent)
		fmt.Printf("%s  Condition:\n", indent)
		PrintAST(n.Condition, indent+"    ")
		
		fmt.Printf("%s  Consequence:\n", indent)
		PrintAST(n.Consequence, indent+"    ")
		
		if n.Alternative != nil {
			fmt.Printf("%s  Alternative:\n", indent)
			PrintAST(n.Alternative, indent+"    ")
		}
	case *FunctionLiteral:
		fmt.Printf("%sFunctionLiteral:\n", indent)
		fmt.Printf("%s  Parameters:\n", indent)
		for _, p := range n.Parameters {
			PrintAST(p, indent+"    ")
		}
		fmt.Printf("%s  Body:\n", indent)
		PrintAST(n.Body, indent+"    ")
	case *BlockStatement:
		fmt.Printf("%sBlockStatement:\n", indent)
		for _, stmt := range n.Statements {
			PrintAST(stmt, indent+"  ")
		}
	case *CallExpression:
		fmt.Printf("%sCallExpression:\n", indent)
		fmt.Printf("%s  Function:\n", indent)
		PrintAST(n.Function, indent+"    ")

		if len(n.Arguments) > 0 {
			fmt.Printf("%s  Arguments:\n", indent)
			for _, arg := range n.Arguments {
				PrintAST(arg, indent+"    ")
			}
		} else {
			fmt.Printf("%s  Arguments: (none)\n", indent)
		}
	default:
		fmt.Printf("%sUnknown Node: %T\n", indent, n)
	}
}