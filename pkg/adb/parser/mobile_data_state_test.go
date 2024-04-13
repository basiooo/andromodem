package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseMobileDataConnectionState(t *testing.T) {
	data := `2`
	expected := parser.DataConnected
	mobileDataConnectionStateInfo := parser.NewMobileDataConnectionState(data)
	actual := mobileDataConnectionStateInfo.ConnectionState
	assert.Equal(t, expected, actual)
}
func TestParseMobileDataConnectionStateInvalidState(t *testing.T) {
	data := `aaa`
	expected := parser.DataUnknown
	mobileDataConnectionStateInfo := parser.NewMobileDataConnectionState(data)
	actual := mobileDataConnectionStateInfo.ConnectionState
	assert.Equal(t, expected, actual)
}
