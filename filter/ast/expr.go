package ast

type Expr interface {
	exprType()
}

type OperationExpr struct {
	Name string
	Expr Expr
}

type AccessorExpr struct {
	Ident string
	Name  string
}

type UnionExpr struct {
	Left  Expr
	Right Expr
}

type AndExpr struct {
	Left Expr
	Right Expr
}

type OrExpr struct {
	Left Expr
	Right Expr
}

func (e *OperationExpr) exprType() {}
func (e *AccessorExpr) exprType()  {}
func (e *UnionExpr) exprType()     {}
func (e *AndExpr) exprType()       {}
func (e *OrExpr) exprType()       {}
