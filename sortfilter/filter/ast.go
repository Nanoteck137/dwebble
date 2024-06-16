package filter

type FilterExpr interface {
	filterExprType()
}

type AndExpr struct {
	Left  FilterExpr
	Right FilterExpr
}

type OrExpr struct {
	Left  FilterExpr
	Right FilterExpr
}

type OpKind int

const (
	OpEqual OpKind = iota
	OpNotEqual
	OpLike
	OpGreater
)

type OpExpr struct {
	Kind  OpKind
	Name  string
	Value any
}

type Table struct {
	Name       string
	SelectName string
	WhereName  string
}

type InTableExpr struct {
	Not   bool
	Table Table
	Ids   []string
}

func (e *AndExpr) filterExprType()     {}
func (e *OrExpr) filterExprType()      {}
func (e *OpExpr) filterExprType()      {}
func (e *InTableExpr) filterExprType() {}
