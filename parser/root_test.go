package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseRoot(t *testing.T) {
	data := "0.7.0:KernelSU"
	expected := parser.Root{
		IsRooted: true,
		RootInfo: &parser.RootInfo{
			Version: "0.7.0",
			Name:    "KernelSU",
		},
	}
	actual := parser.ParseRoot(data)
	assert.Equal(t, expected, actual)
}
func TestParseRootNotRooted(t *testing.T) {
	data := "/system/bin/sh: su: not found"
	expected := parser.Root{
		IsRooted: false,
	}
	actual := parser.ParseRoot(data)
	assert.Equal(t, expected, actual)
}
