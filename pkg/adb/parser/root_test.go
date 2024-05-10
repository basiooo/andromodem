package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseRoot(t *testing.T) {
	data := "0.7.0:KernelSU"
	expected := parser.Root{
		IsRooted: true,
		RootDetail: &parser.RootDetail{
			Version: "0.7.0",
			Name:    "KernelSU",
		},
	}
	rootInfo := parser.NewRoot(data)
	actual := *rootInfo
	assert.Equal(t, expected, actual)
}
func TestParseRootOldVersion(t *testing.T) {
	data := "16 superuser"
	expected := parser.Root{
		IsRooted: true,
		RootDetail: &parser.RootDetail{
			Version: "16",
			Name:    "superuser",
		},
	}
	rootInfo := parser.NewRoot(data)
	actual := *rootInfo
	assert.Equal(t, expected, actual)
}
func TestParseRootNotRooted(t *testing.T) {
	data := "/system/bin/sh: su: not found"
	expected := parser.Root{
		IsRooted: false,
	}
	rootInfo := parser.NewRoot(data)
	actual := *rootInfo
	assert.Equal(t, expected, actual)
}
