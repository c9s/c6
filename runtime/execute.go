package runtime

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

const MaxWhileIterations = 10_000

type Printer func(msg any)

type Runtime struct {
	DebugPrinter Printer
	WarnPrinter  Printer
}

func NewRuntime(debug, warn Printer) *Runtime {
	return &Runtime{
		DebugPrinter: debug,
		WarnPrinter:  warn,
	}
}

func (r *Runtime) ExecuteList(scope *Scope, stmts *ast.StmtList) (*ast.StmtList, error) {
	out := &ast.StmtList{}

	for _, stmt := range stmts.Stmts {
		ret, err := r.ExecuteSingle(scope, stmt)

		if err != nil {
			return nil, err
		}

		out.AppendList(ret)
	}

	return out, nil
}

func (r *Runtime) ExecuteSingle(scope *Scope, stmt ast.Stmt) (*ast.StmtList, error) {
	switch t := stmt.(type) {
	case *ast.LogStmt:
		err := r.executeLogStmt(scope, t)
		return nil, err
	case *ast.AssignStmt:
		err := r.executeAssignStmt(scope, t)
		return nil, err
	case *ast.IfStmt:
		return r.executeIfStmt(scope, t)
	case *ast.ForStmt:
		return r.executeForStmt(scope, t)
	case *ast.WhileStmt:
		return r.executeWhileStmt(scope, t)
	case *ast.RuleSet:
		return r.executeRuleSet(scope, t)
	case *ast.Property:
		return r.executeProperty(scope, t)
	}

	return nil, fmt.Errorf("Don't know how to execute the statement %v", stmt)
}

func (r *Runtime) executeCondition(e ast.Expr, scope *Scope) (bool, error) {
	v, err := EvaluateExprInBooleanContext(e, scope)

	if err != nil {
		return false, err
	}

	if v == nil {
		return false, nil
	}

	if bval, ok := v.(ast.BooleanValue); ok {
		return bval.Boolean(), nil
	}

	if _, ok := v.(ast.Null); ok {
		return false, nil
	}

	return false, fmt.Errorf("BooleanValue interface is not support for %T", v)
}

func (r *Runtime) executeWhileStmt(scope *Scope, stmt *ast.WhileStmt) (*ast.StmtList, error) {
	child := NewScope(scope)
	count := 0
	out := &ast.StmtList{}

	for {
		keepGoing, err := r.executeCondition(stmt.Condition, child)

		if err != nil {
			return nil, err
		}

		if !keepGoing {
			break
		}

		l, err := r.ExecuteList(child, &stmt.Block.Stmts)

		if err != nil {
			return nil, err
		}

		out.AppendList(l)
		count++

		if count == MaxWhileIterations {
			return nil, fmt.Errorf("While loop had %d iterations, most probably it's an infinite loop", count)
		}
	}

	return out, nil
}

func (r *Runtime) executeForStmt(scope *Scope, stmt *ast.ForStmt) (*ast.StmtList, error) {

	eval := func(e ast.Expr) (int, error) {
		v, err := EvaluateExpr(e, scope)

		if err != nil {
			return 0, err
		}

		num, ok := v.(*ast.Number)

		if !ok {
			return 0, fmt.Errorf("Expected to get a number but got %T", v)
		}

		return num.Integer(), nil
	}

	from, err := eval(stmt.From)

	if err != nil {
		return nil, err
	}

	to, err := eval(stmt.To)

	if err != nil {
		return nil, err
	}

	if from == to {
		return nil, nil
	}

	step := 1

	if to < from {
		step = -1
	}

	// to simplify the loop,
	// since we can now treat both loops as exclusive upper limit
	if stmt.Inclusive {
		to = to + step
	}

	out := &ast.StmtList{}

	for from != to {
		child := NewScope(scope)
		child.Insert(stmt.Variable.NormalizedName(), ast.NewNumber(float64(from), nil, nil))

		l, err := r.ExecuteList(child, &stmt.Block.Stmts)

		if err != nil {
			return nil, err
		}

		out.AppendList(l)
		from = from + step
	}

	return out, nil
}

func (r *Runtime) executeIfStmt(scope *Scope, stmt *ast.IfStmt) (*ast.StmtList, error) {
	var out *ast.StmtList

	v, err := r.executeCondition(stmt.Condition, scope)

	if err != nil {
		return nil, err
	}

	if v {
		out = &stmt.Block.Stmts
	} else {
		for _, elsif := range stmt.ElseIfs {
			v, err := r.executeCondition(elsif.Condition, scope)

			if err != nil {
				return nil, err
			}

			if v {
				out = &elsif.Block.Stmts
				break
			}
		}

		if out == nil && stmt.ElseBlock != nil {
			out = &stmt.ElseBlock.Stmts
		}
	}

	if out == nil {
		return nil, nil
	}

	child := NewScope(scope)

	return r.ExecuteList(child, out)
}

func (r *Runtime) executeLogStmt(scope *Scope, stmt *ast.LogStmt) error {
	if stmt.Expr == nil {
		return fmt.Errorf("Expected expression")
	}

	v, err := EvaluateExpr(stmt.Expr, scope)

	if err != nil {
		return err
	}

	switch stmt.LogLevel {
	case ast.LogLevelDebug:
		r.DebugPrinter(v)
	case ast.LogLevelWarn:
		r.WarnPrinter(v)
	default:
		return fmt.Errorf(v.String())
	}

	return nil
}

func (r *Runtime) executeAssignStmt(scope *Scope, stmt *ast.AssignStmt) error {
	varName := stmt.Variable.Name
	val, err := EvaluateExpr(stmt.Expr, scope)

	if err != nil {
		return err
	}

	if stmt.Global {
		scope.GetGlobal().Insert(varName, val)
	} else {
		scope.Insert(varName, val)
	}

	return nil
}

func (r *Runtime) executeRuleSet(scope *Scope, stmt *ast.RuleSet) (*ast.StmtList, error) {
	child := NewScope(scope)

	res, err := r.ExecuteList(child, &stmt.Block.Stmts)

	if err != nil {
		return nil, err
	}

	rs := ast.NewRuleSet()
	decl := ast.NewDeclBlock(rs)
	decl.AppendList(res)
	rs.Block = decl
	rs.Selectors = stmt.Selectors

	return &ast.StmtList{
		Stmts: []ast.Stmt{rs},
	}, nil
}

func (r *Runtime) executeProperty(scope *Scope, stmt *ast.Property) (*ast.StmtList, error) {
	ret := ast.NewProperty(stmt.Name.Token)

	for _, e := range stmt.Values {
		val, err := EvaluateExpr(e, scope)
		if err != nil {
			return nil, err
		}

		ret.Values = append(ret.Values, val)
	}

	return &ast.StmtList{
		Stmts: []ast.Stmt{ret},
	}, nil
}
