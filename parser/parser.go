package parser

import (
	"interpreter_go/ast"
	"interpreter_go/lexer"
	"interpreter_go/token"
	"strconv"
	"fmt"
)

const (
	_ int = iota
	LOWEST
	EQUALS // ==
	LESSGREATER // > or <
	SUM // +
	PRODUCT
	PREFIX
	CALL
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn
	errors []string
}

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		l: lex,
		errors: []string{},
	}
	
	parser.nextToken()
	parser.nextToken()

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)

	return parser
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
		case token.LET:
			return p.parseLetStatement()
		case token.RETURN:
			return p.parseReturnStatement()
		case token.IDENT:
			return p.parseExpressionStatement()
		case token.INT:
			return p.parseExpressionStatement()
		default:
			return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	litt := &ast.IntegerLiteral{Token: p.currToken}

	val, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	litt.Value = val
	return litt
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
} 

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}


func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
} 

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	error := fmt.Sprintf("Expected token: %s got %s\n", t, p.peekToken.Type)
	p.errors = append(p.errors, error)
}
