package loglint

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer reports common issues in log messages: case, language, special chars, sensitive data.
var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "checks log messages for style and security issues",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		msgExpr, ok := extractLogMessage(pass, call)
		if !ok {
			return
		}

		checkAllRules(pass, msgExpr)
	})

	return nil, nil
}

// extractLogMessage returns the message argument expression if the call is a
// recognized logging call from log/slog or go.uber.org/zap.
func extractLogMessage(pass *analysis.Pass, call *ast.CallExpr) (ast.Expr, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	methodName := sel.Sel.Name
	msgIdx, known := logMethodMsgIndex(methodName)
	if !known {
		return nil, false
	}

	if !isFromSupportedLogger(pass, sel) {
		return nil, false
	}

	if len(call.Args) <= msgIdx {
		return nil, false
	}
	return call.Args[msgIdx], true
}

func logMethodMsgIndex(name string) (int, bool) {
	switch name {
	// standard log methods (slog + zap)
	case "Info", "Warn", "Error", "Debug":
		return 0, true
	// zap-only levels
	case "Fatal", "Panic", "DPanic":
		return 0, true
	// zap sugared logger
	case "Infow", "Warnw", "Errorw", "Debugw",
		"Fatalw", "Panicw", "DPanicw",
		"Infof", "Warnf", "Errorf", "Debugf",
		"Fatalf", "Panicf", "DPanicf":
		return 0, true
	// slog context variants — message is second arg
	case "InfoContext", "WarnContext", "ErrorContext", "DebugContext":
		return 1, true
	// slog.Log / LogAttrs — ctx, level, msg
	case "Log", "LogAttrs":
		return 2, true
	}
	return -1, false
}

func isFromSupportedLogger(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	// package-level call: slog.Info(...)
	if ident, ok := sel.X.(*ast.Ident); ok {
		obj := pass.TypesInfo.Uses[ident]
		if obj == nil {
			return false
		}
		if pkgName, ok := obj.(*types.PkgName); ok {
			p := pkgName.Imported().Path()
			return p == "log/slog" || p == "go.uber.org/zap"
		}
	}

	// method call on *slog.Logger / *zap.Logger / *zap.SugaredLogger
	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ == nil {
		return false
	}
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj.Pkg() == nil {
		return false
	}

	pkgPath := obj.Pkg().Path()
	typeName := obj.Name()

	switch {
	case pkgPath == "log/slog" && typeName == "Logger":
		return true
	case pkgPath == "go.uber.org/zap" && (typeName == "Logger" || typeName == "SugaredLogger"):
		return true
	}
	return false
}

// stringLiteral pairs the unquoted value with its AST node for position info.
type stringLiteral struct {
	value string
	node  *ast.BasicLit
}

// extractStringLiterals collects every string literal inside expr,
// recursing into binary "+" concatenations and call arguments (e.g. fmt.Sprintf).
func extractStringLiterals(expr ast.Expr) []stringLiteral {
	var result []stringLiteral
	ast.Inspect(expr, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING {
			val, err := strconv.Unquote(lit.Value)
			if err == nil {
				result = append(result, stringLiteral{value: val, node: lit})
			}
		}
		return true
	})
	return result
}

// hasNonLiteralParts reports whether expr contains anything besides string constants.
func hasNonLiteralParts(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return e.Kind != token.STRING
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return hasNonLiteralParts(e.X) || hasNonLiteralParts(e.Y)
		}
		return true
	case *ast.ParenExpr:
		return hasNonLiteralParts(e.X)
	default:
		return true
	}
}

func fullMessage(lits []stringLiteral) string {
	var b strings.Builder
	for _, l := range lits {
		b.WriteString(l.value)
	}
	return b.String()
}
