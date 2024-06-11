package ast

import "github.com/nanoteck137/dwebble/filter/token"

type Expr interface {
	exprType()
}

type CallExpr struct {
	Name   string
	Params []Expr
}

type FieldExpr struct {
	Ident string
	Field string
}

type OpExpr struct {
	Kind  token.Kind
	Name  string
	Value any
}

type AndExpr struct {
	Left  Expr
	Right Expr
}

type OrExpr struct {
	Left  Expr
	Right Expr
}

type Table struct {
	Type string
	Ids  []string
}

type InTableExpr struct {
	Not   bool
	Table Table
}

func (e *CallExpr) exprType()    {}
func (e *FieldExpr) exprType()   {}
func (e *OpExpr) exprType()      {}
func (e *AndExpr) exprType()     {}
func (e *OrExpr) exprType()      {}
func (e *InTableExpr) exprType() {}
