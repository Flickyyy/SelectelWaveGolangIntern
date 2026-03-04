package main

import (
	loglint "github.com/Flickyyy/SelectelWaveGolangIntern"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loglint.Analyzer)
}
