package evaluator

import (
	"interpreter_go/ast"
	"interpreter_go/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
		case *ast.Program:
			return evalStatements(node.Statements)
		case *ast.ExpressionStatement:
			return Eval(node.Expression)
		case *ast.IntegerLiteral:
			return &object.Integer{Value: node.Value}
	}
	return nil
}

func evalStatements(statements []ast.Statement) object.Object {

	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}	