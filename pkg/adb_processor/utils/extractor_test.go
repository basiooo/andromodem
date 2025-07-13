package utils_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/basiooo/andromodem/pkg/adb_processor/utils"
	"github.com/stretchr/testify/assert"
)

type FakeParser struct{}

func (f *FakeParser) Parse(output string) error {
	return nil
}

func TestGetResultFromRaw_Valid(t *testing.T) {
	t.Parallel()
	expected := "device_model_123"
	raw := &parser.RawParser{
		Result: expected,
	}
	result := utils.GetResultFromRaw(raw)
	assert.Equal(t, expected, result)
}

func TestGetResultFromRaw_InvalidType(t *testing.T) {
	t.Parallel()

	fake := &FakeParser{}
	result := utils.GetResultFromRaw(fake)
	assert.Equal(t, "", result)
}

func TestGetResultFromRaw_Nil(t *testing.T) {
	var result parser.IParser = nil
	output := utils.GetResultFromRaw(result)
	assert.Equal(t, "", output)
}
