package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseBussyboxInstalled(t *testing.T) {
	t.Parallel()
	data := `/system/xbin/busybox`
	expected := true
	busyboxCheck := parser.NewBusyboxCheck()
	err := busyboxCheck.Parse(data)
	assert.NoError(t, err)
	result := busyboxCheck.(*parser.BusyboxCheck)
	assert.Equal(t, expected, result.Exist)
}
func TestParseBussyboxNotInstalled(t *testing.T) {
	t.Parallel()
	data := ``
	expected := false
	busyboxCheck := parser.NewBusyboxCheck()
	err := busyboxCheck.Parse(data)
	assert.NoError(t, err)
	result := busyboxCheck.(*parser.BusyboxCheck)
	assert.Equal(t, expected, result.Exist)
}

func BenchmarkParseBussyboxInstalled(b *testing.B) {
	data := `/system/xbin/busybox`
	for i := 0; i < b.N; i++ {
		busyboxCheck := parser.NewBusyboxCheck()
		_ = busyboxCheck.Parse(data)
	}
}

func BenchmarkParseBussyboxNotInstalled(b *testing.B) {
	data := ``
	for i := 0; i < b.N; i++ {
		busyboxCheck := parser.NewBusyboxCheck()
		_ = busyboxCheck.Parse(data)
	}
}
