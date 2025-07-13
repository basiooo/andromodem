package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseStorage(t *testing.T) {
	t.Parallel()
	data := "Latency: 0ms [512B Data Write]\nData-Free: 5684936K / 6094400K total = 93% free\nCache-Free: 5684936K / 6094400K total = 93% free\nSystem-Free: 1098536K / 2539312K total = 43% free"
	expected := &parser.Storage{
		SystemTotal: 2539312,
		SystemFree:  1098536,
		SystemUsed:  1440776,
		DataTotal:   6094400,
		DataFree:    5684936,
		DataUsed:    409464,
	}
	storage := parser.NewStorage()
	err := storage.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, storage)
}

func TestParseStorageInvalidDataFree(t *testing.T) {
	t.Parallel()
	data := "Latency: 0ms [512B Data Write]\nCache-Free: 5684936K / 6094400K total = 93% free\nSystem-Free: 1098536K / 2539312K total = 43% free"
	// expected := &parser.Storage{
	// 	SystemTotal: 2539312,
	// 	SystemFree:  1098536,
	// 	SystemUsed:  1440776,
	// }
	storage := parser.NewStorage()
	err := storage.Parse(data)
	assert.NoError(t, err)
	// assert.Equal(t, expected, storage)
}

func BenchmarkParseStorage(b *testing.B) {
	data := "Latency: 0ms [512B Data Write]\nData-Free: 5684936K / 6094400K total = 93% free\nCache-Free: 5684936K / 6094400K total = 93% free\nSystem-Free: 1098536K / 2539312K total = 43% free"
	for i := 0; i < b.N; i++ {
		storage := parser.NewStorage()
		_ = storage.Parse(data)
	}
}

func BenchmarkParseStorageInvalidDataFree(b *testing.B) {
	data := "Latency: 0ms [512B Data Write]\nCache-Free: 5684936K / 6094400K total = 93% free\nSystem-Free: 1098536K / 2539312K total = 43% free"
	for i := 0; i < b.N; i++ {
		storage := parser.NewStorage()
		_ = storage.Parse(data)
	}
}
