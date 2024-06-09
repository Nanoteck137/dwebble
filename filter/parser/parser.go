package parser

import (
	"fmt"
	"io"

	"github.com/nanoteck137/dwebble/filter/ast"
	"github.com/nanoteck137/dwebble/filter/lexer"
	"github.com/nanoteck137/dwebble/filter/token"
)

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

type ErrorList []*Error

func (p *ErrorList) Add(message string) {
	*p = append(*p, &Error{
		Message: message,
	})
}

type Parser struct {
	tokenizer *lexer.Tokenizer
	token     token.Token
	Errors    ErrorList
}

func New(reader io.Reader) *Parser {
	tokenizer := lexer.New(reader)

	return &Parser{
		tokenizer: tokenizer,
		token:     tokenizer.NextToken(),
	}
}

func (p *Parser) error(message string) {
	p.Errors.Add(message)
}

func (p *Parser) next() {
	p.token = p.tokenizer.NextToken()
}

func (p *Parser) expect(token token.Kind) {
	if p.token.Kind != token {
		p.error(fmt.Sprintf("Expected token: %v got %v", token, p.token.Kind))
	}

	p.next()
}

func (p *Parser) parseExprBase() ast.Expr {
	if p.token.Kind == token.Ident {
		name := p.token.Ident
		p.next()

		p.expect(token.Dot)

		ident := p.token.Ident
		p.expect(token.Str)

		return &ast.AccessorExpr{
			Ident: ident,
			Name:  name,
		}
	}

	return nil
}

func (p *Parser) parseExpr0() ast.Expr {
	if p.token.Kind == token.Ident {
		name := p.token.Ident
		p.next()

		p.expect(token.LParen)

		e := p.parseExprBase()

		p.expect(token.RParen)

		return &ast.OperationExpr{
			Name: name,
			Expr: e,
		}
	}

	return nil
}

func (p *Parser) parseExpr1() ast.Expr {
	left := p.parseExpr0()

	for p.token.Kind == token.DoubleAnd || p.token.Kind == token.DoubleOr {
		op := p.token.Kind
		p.next()

		right := p.parseExpr0()

		switch op {
		case token.DoubleAnd:
			left = &ast.AndExpr{
				Left:  left,
				Right: right,
			}
		case token.DoubleOr:
			left = &ast.OrExpr{
				Left:  left,
				Right: right,
			}
		default:
			panic("Unhandled case")
		}
	}

	return left
}

func (p *Parser) ParseExpr() ast.Expr {
	return p.parseExpr1()
}
