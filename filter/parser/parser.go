package parser

import (
	"fmt"
	"io"

	"github.com/nanoteck137/dwebble/filter/ast"
	"github.com/nanoteck137/dwebble/filter/lexer"
	"github.com/nanoteck137/dwebble/filter/token"
	"github.com/nanoteck137/dwebble/types"
)

type Parser struct {
	tokenizer *lexer.Tokenizer
	token     token.Token
	Errors    types.ErrorList
}

func New(reader io.Reader) *Parser {
	tokenizer := lexer.New(reader)

	p := &Parser{
		tokenizer: tokenizer,
		token:     tokenizer.NextToken(),
		Errors:    []*types.Error{},
	}

	p.tokenizer.Init(p.error)

	return p
}

func (p *Parser) error(message string) {
	p.Errors.Add(message)
}

func (p *Parser) next() {
	p.token = p.tokenizer.NextToken()
}

func (p *Parser) is(tok token.Kind) bool {
	return p.token.Kind == tok
}

func (p *Parser) expect(token token.Kind) {
	if p.token.Kind != token {
		p.error(fmt.Sprintf("Expected token: %v got %v", token.String(), p.token.Kind.String()))
	}

	p.next()
}

func (p *Parser) parseCallParam() ast.Expr {
	if p.token.Kind == token.Ident {
		ident := p.token.Ident
		p.next()

		p.expect(token.Dot)

		field := p.token.Ident
		p.expect(token.Str)

		return &ast.FieldExpr{
			Ident: ident,
			Field: field,
		}
	}

	return nil
}

func (p *Parser) parseExpr0() ast.Expr {
	switch p.token.Kind {
	case token.Ident:
		name := p.token.Ident
		p.next()

		switch p.token.Kind {
		case token.LParen:
			p.next()

			var params []ast.Expr
			if p.token.Kind != token.RParen {
				e := p.parseCallParam()
				params = append(params, e)

				for p.token.Kind != token.RParen {
					p.expect(token.Comma)
					e := p.parseCallParam()

					params = append(params, e)
				}
			}

			p.expect(token.RParen)

			return &ast.CallExpr{
				Name:   name,
				Params: params,
			}
		case token.DoubleEqual:
			kind := p.token.Kind
			p.next()

			value := p.token.Ident
			p.expect(token.Str)

			return &ast.OpExpr{
				Kind:  kind,
				Name:  name,
				Value: value,
			}
		}

	case token.LParen:
		p.next()

		e := p.ParseExpr()
		p.expect(token.RParen)

		return e
	}

	p.error("Unexpected token: " + p.token.Kind.String())

	return nil
}

func (p *Parser) parseExpr1() ast.Expr {
	left := p.parseExpr0()

	for p.is(token.DoubleAnd) || p.is(token.DoubleOr) {
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

func (p *Parser) Parse() ast.Expr {
	e := p.ParseExpr()
	p.expect(token.Eof)

	return e
}
