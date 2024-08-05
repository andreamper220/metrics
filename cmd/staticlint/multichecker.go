package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"regexp"
)

var osExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit() in main",
	Run:  runOsExitCheck,
}

func main() {
	// add standard analysers + osexitcheck
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		osExitCheckAnalyzer,
	}

	// add SA* analysers
	sapattern := regexp.MustCompile(`SA\d+`)
	for _, v := range staticcheck.Analyzers {
		if sapattern.MatchString(v.Analyzer.Name) {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// add ST1019 analyser
	for _, v := range stylecheck.Analyzers {
		if v.Analyzer.Name == "ST1019" {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// add QF1002 analyser
	for _, v := range quickfix.Analyzers {
		if v.Analyzer.Name == "QF1002" {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}

func runOsExitCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if funcDecl, ok := node.(*ast.FuncDecl); ok {
				if funcDecl.Name.String() == "main" {
					bodyList := funcDecl.Body.List
					for _, stmt := range bodyList {
						if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
							if call, ok := exprStmt.X.(*ast.CallExpr); ok {
								if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
									xName := fun.X.(*ast.Ident).Name
									funcName := fun.Sel.Name
									if xName == "os" && funcName == "Exit" {
										pass.Reportf(fun.Pos(), "os.Exit() within main function")
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
