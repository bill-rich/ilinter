package ilinter_test

import (
	"testing"

	"github.com/bill-rich/ilinter"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ilinter.Analyzer, "a")
}
