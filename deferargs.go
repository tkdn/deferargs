package deferargs

import (
	"go/ast"

	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "deferargs",
	Doc:      "reports defer calls whose arguments are variable references (immediately evaluated)",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Preorder([]ast.Node{(*ast.DeferStmt)(nil)}, func(n ast.Node) {
		ds := n.(*ast.DeferStmt)

		call := ds.Call
		// function literal = クロージャの場合は処理を終える
		if _, isLit := call.Fun.(*ast.FuncLit); isLit {
			return
		}

		hasVarArg := false
		for _, arg := range call.Args {
			if isVarRef(pass.TypesInfo, arg) {
				hasVarArg = true
				break
			}
		}

		if hasVarArg {
			pass.Reportf(call.Lparen, "defer with variable argument(s): value(s) are evaluated immediately; wrap in a closure if a later update is intended")
		}
	})
	return nil, nil
}

func isVarRef(info *types.Info, expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		_, ok := info.Uses[e].(*types.Var)
		return ok
	case *ast.SelectorExpr:
		// obj.Field などを確認。Sel が *types.Var ならフィールドかpkg変数
		_, ok := info.Uses[e.Sel].(*types.Var)
		return ok
	default:
		return false
	}
}
