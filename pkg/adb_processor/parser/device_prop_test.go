package parser_test

import (
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDeviceProp(t *testing.T) {
	t.Parallel()
	data := "[ro.product.model]: [Phone]\n[ro.product.brand]: [Custom]\n[ro.product.name]: [vbox86p]\n,[ro.product.name]: [vbox86p]\n"
	expected := &parser.DeviceProp{
		Model: "Phone",
		Brand: "Custom",
		Name:  "vbox86p",
	}
	deviceProp := parser.NewDeviceProp()
	err := deviceProp.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, deviceProp)
}

func TestParseDevicePropHasEmptyValue(t *testing.T) {
	t.Parallel()
	data := "[ro.product.model]: \n[ro.product.brand]: [Custom]\n[ro.product.name]: [vbox86p]\n,[ro.product.name]: [vbox86p]\n"
	expected := &parser.DeviceProp{
		Model: "",
		Brand: "Custom",
		Name:  "vbox86p",
	}
	deviceProp := parser.NewDeviceProp()
	err := deviceProp.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, deviceProp)
}

func BenchmarkParseDeviceProp(b *testing.B) {
	data := "[ro.product.model]: [Phone]\n[ro.product.brand]: [Custom]\n[ro.product.name]: [vbox86p]\n,[ro.product.name]: [vbox86p]\n"
	for i := 0; i < b.N; i++ {
		deviceProp := parser.NewDeviceProp()
		_ = deviceProp.Parse(data)
	}
}

func BenchmarkParseDevicePropHasEmptyValue(b *testing.B) {
	data := "[ro.product.model]: \n[ro.product.brand]: [Custom]\n[ro.product.name]: [vbox86p]\n,[ro.product.name]: [vbox86p]\n"
	for i := 0; i < b.N; i++ {
		deviceProp := parser.NewDeviceProp()
		_ = deviceProp.Parse(data)
	}
}
