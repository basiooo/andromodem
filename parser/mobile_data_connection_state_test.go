package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseMobileDataConnectionState(t *testing.T) {
	data := `-1
	2`
	expected := []parser.MobileDataConnectionState{parser.DataUnknown, parser.DataConnected}
	actual := parser.ParseMobileDataConnectionState(data)
	assert.Equal(t, expected, actual)
}
func TestParseMobileDataConnectionStateInvalidState(t *testing.T) {
	data := `-aaa
	test`
	expected := []parser.MobileDataConnectionState{parser.DataUnknown, parser.DataUnknown}
	actual := parser.ParseMobileDataConnectionState(data)
	assert.Equal(t, expected, actual)
}
