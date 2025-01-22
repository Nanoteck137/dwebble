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

type IsNullExpr struct {
	Name string
	Not  bool
}

type OpKind int

const (
	OpEqual OpKind = iota
	OpNotEqual
	OpLike
	OpGreater
	OpGreaterEqual
	OpLesser
	OpLesserEqual
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
	Not        bool
	IdSelector string
	Table      Table
	Ids        []string
}

func (e *AndExpr) filterExprType()     {}
func (e *OrExpr) filterExprType()      {}
func (e *IsNullExpr) filterExprType()  {}
func (e *OpExpr) filterExprType()      {}
func (e *InTableExpr) filterExprType() {}
