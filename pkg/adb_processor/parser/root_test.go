package parser_test

import (
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseRoot(t *testing.T) {
	t.Parallel()
	data := "0.7.0:KernelSU"
	expected := &parser.Root{
		IsRooted: true,
		RootDetail: parser.RootDetail{
			Version: "0.7.0",
			Name:    "KernelSU",
		},
	}
	rootInfo := parser.NewRoot()
	err := rootInfo.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, rootInfo)
}

func TestParseRootOlder(t *testing.T) {
	t.Parallel()
	data := "16 superuser"
	expected := &parser.Root{
		IsRooted: true,
		RootDetail: parser.RootDetail{
			Version: "16",
			Name:    "superuser",
		},
	}
	rootInfo := parser.NewRoot()
	err := rootInfo.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, rootInfo)
}

func TestParseRootNotRooted(t *testing.T) {
	t.Parallel()
	data := "/system/bin/sh: su: not found"
	expected := &parser.Root{
		IsRooted: false,
	}
	rootInfo := parser.NewRoot()
	err := rootInfo.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, rootInfo)
}

func TestParseRootEmpty(t *testing.T) {
	t.Parallel()
	data := ""
	expected := &parser.Root{
		IsRooted: false,
	}
	rootInfo := parser.NewRoot()
	err := rootInfo.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, rootInfo)
}

func TestParseRootUnformatted(t *testing.T) {
	t.Parallel()
	data := "Lorem ipsum dolor sit amet"
	expected := &parser.Root{
		IsRooted: true,
		RootDetail: parser.RootDetail{
			Name: data,
		},
	}
	rootInfo := parser.NewRoot()
	err := rootInfo.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, rootInfo)
}

func BenchmarkParseRoot(b *testing.B) {
	data := "0.7.0:KernelSU"
	for i := 0; i < b.N; i++ {
		rootInfo := parser.NewRoot()
		_ = rootInfo.Parse(data)
	}
}
