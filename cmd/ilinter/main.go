// Command ilinter reports variable and constant names shorter than a minimum
// length. Run it directly on packages or through go vet:
//
//	ilinter ./...
//	go vet -vettool=$(which ilinter) ./...
package main

import (
	"github.com/bill-rich/ilinter"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(ilinter.Analyzer) }
