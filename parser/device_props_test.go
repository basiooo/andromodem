package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseDeviceProps(t *testing.T) {
	data := `[ro.product.model]: [Samsung J2 Prime]\n[ro.product.brand]: [Samsung]\n[ro.product.name]: [SM-G532G]\n[ro.build.version.release]: [14]`
	expected := parser.DeviceProps{
		Model:          "Samsung J2 Prime",
		Brand:          "Samsung",
		Name:           "SM-G532G",
		AndroidVersion: "14",
	}
	actual := parser.ParseDeviceProps(data)
	assert.Equal(t, expected, actual)
}
