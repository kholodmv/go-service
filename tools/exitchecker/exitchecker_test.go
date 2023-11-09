package exitchecker_test

import (
	"github.com/kholodmv/go-service/tools/exitchecker"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMyAnalyzer(t *testing.T) {
	// function exitchecker.run applies the analyzed analyzer ErrCheckAnalyzer
	// to packages from the testdata folder and checks expectations
	// ./... - checking all subdirectories in testdata
	// you can specify ./test to check only test
	analysistest.Run(t, analysistest.TestData(), exitchecker.Analyzer, "./...")
}
