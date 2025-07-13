package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
)

func BenchmarkParseRaw(b *testing.B) {
	data := "This is a raw string that needs to be parsed"
	for i := 0; i < b.N; i++ {
		rawParser := parser.NewRawParser()
		_ = rawParser.Parse(data)
	}
}