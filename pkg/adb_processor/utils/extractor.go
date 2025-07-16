package utils

import "github.com/basiooo/andromodem/pkg/adb_processor/parser"

func GetResultFromRaw(result parser.IParser) string {
	if parsed, ok := result.(*parser.RawParser); ok {
		return parsed.Result
	}
	return ""
}
