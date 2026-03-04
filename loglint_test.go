package loglint_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	loglint "github.com/Flickyyy/SelectelWaveGolangIntern"
)

func TestSlogRules(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, loglint.Analyzer, "slogcheck")
}

func TestZapRules(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, loglint.Analyzer, "zapcheck")
}
