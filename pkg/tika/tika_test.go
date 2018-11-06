package tika

import (
	"path/filepath"
	"testing"
)

func TestScraper(t *testing.T) {
	absPath, _ := filepath.Abs(".")
	testFilePath := filepath.Join(absPath, "testdata", "b9c2c0469b9101d66667b79166061f7e4d62c9ef")
	GetTikaContentAndMetaForFile(testFilePath)
}
