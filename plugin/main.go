//go:build ignore

package main

import (
	loglint "github.com/Flickyyy/SelectelWaveGolangIntern"
	"golang.org/x/tools/go/analysis"
)

// AnalyzerPlugin implements the interface golangci-lint expects for Go plugins.
type AnalyzerPlugin struct{}

func (*AnalyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{loglint.Analyzer}
}

// Plugin is the exported symbol golangci-lint looks for.
var Plugin AnalyzerPlugin
