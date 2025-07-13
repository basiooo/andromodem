package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseMobileDataState(t *testing.T) {
	t.Parallel()
	data := `2`
	expected := parser.DataConnected
	mobileDataState := parser.NewMobileDataState()
	err := mobileDataState.Parse(data)
	assert.NoError(t, err)
	state := mobileDataState.(*parser.MobileDataState)
	assert.Equal(t, expected, state.State)
}
func TestParseMobileDataStateToString(t *testing.T) {
	t.Parallel()
	data := `2`
	expected := "Connected"
	mobileDataState := parser.NewMobileDataState()
	err := mobileDataState.Parse(data)
	assert.NoError(t, err)
	state := mobileDataState.(*parser.MobileDataState)
	assert.Equal(t, expected, state.String())
}
func TestParseMobileDataStateInvalidState(t *testing.T) {
	t.Parallel()
	data := `invalid123`
	expected := parser.DataUnknown
	mobileDataState := parser.NewMobileDataState()
	err := mobileDataState.Parse(data)
	assert.NoError(t, err)
	state := mobileDataState.(*parser.MobileDataState)
	assert.Equal(t, expected, state.State)
}

func BenchmarkParseMobileDataState(b *testing.B) {
	data := `2`
	for i := 0; i < b.N; i++ {
		mobileDataState := parser.NewMobileDataState()
		_ = mobileDataState.Parse(data)
	}
}

func BenchmarkParseMobileDataStateInvalidState(b *testing.B) {
	data := `invalid123`
	for i := 0; i < b.N; i++ {
		mobileDataState := parser.NewMobileDataState()
		_ = mobileDataState.Parse(data)
	}
}
