package main

import (
	"github.com/tkdn/deferargs"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(deferargs.Analyzer) }
