package parser_test

import (
	"github.com/basiooo/andromodem/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseThermal(t *testing.T) {
	data := `pmi632-tz
37000
pmi632-ibat-lvl0
-234
			`
	expected := parser.Thermal{
		Data: map[string]int{
			"pmi632-ibat-lvl0": -234,
			"pmi632-tz":        37,
		},
	}
	actual := parser.ParseThermal(data)
	assert.Equal(t, expected, actual)
}
