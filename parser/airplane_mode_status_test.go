package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/parser"
	"github.com/stretchr/testify/assert"
)

func TestAirplaneModeStatus(t *testing.T) {
	data := `enabled`
	actual := parser.ParseAirplaneModeStatus(data)
	assert.Equal(t, true, actual)
}
func TestAirplaneModeStatusInvalid(t *testing.T) {
	data := `test`
	actual := parser.ParseAirplaneModeStatus(data)
	assert.Equal(t, false, actual)
}
