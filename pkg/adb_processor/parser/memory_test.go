package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseMemory(t *testing.T) {
	t.Parallel()
	data := "MemTotal:        1982896 kB\nMemFree:          184732 kB\nMemAvailable:    1258100 kB\nBuffers:           30676 kB\nCached:          1103848 kB\nSwapCached:            0 kB\nActive:          1126396 kB\nInactive:         475672 kB\nActive(anon):     470064 kB\nInactive(anon):      600 kB\nActive(file):     656332 kB\nInactive(file):   475072 kB\nUnevictable:        2844 kB\nMlocked:            2844 kB\nSwapTotal:             0 kB\nSwapFree:              0 kB\nDirty:                 0 kB\nWriteback:             0 kB\nAnonPages:        470460 kB\nMapped:           533244 kB\nShmem:               900 kB\nSlab:              81524 kB\nSReclaimable:      31544 kB\nSUnreclaim:        49980 kB\nKernelStack:       16016 kB\nPageTables:        24416 kB\nNFS_Unstable:          0 kB\nBounce:                0 kB\nWritebackTmp:          0 kB\nCommitLimit:      991448 kB\nCommitted_AS:   33535472 kB\nVmallocTotal:   258867136 kB\nVmallocUsed:           0 kB\nVmallocChunk:          0 kB\nCmaTotal:          40960 kB\nCmaFree:            2280 kB\nHugePages_Total:       0\nHugePages_Free:        0\nHugePages_Rsvd:        0\nHugePages_Surp:        0\nHugepagesize:       2048 kB"
	expected := &parser.Memory{
		MemTotal: 1982896,
		MemFree:  1258100,
		MemUsed:  724796,
	}
	memory := parser.NewMemory()
	err := memory.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, memory)
}
func TestParseMemoryWithSwap(t *testing.T) {
	t.Parallel()
	data := "MemTotal:        1982896 kB\nMemFree:          184732 kB\nMemAvailable:    1258100 kB\nBuffers:           30676 kB\nCached:          1103848 kB\nSwapCached:            0 kB\nActive:          1126396 kB\nInactive:         475672 kB\nActive(anon):     470064 kB\nInactive(anon):      600 kB\nActive(file):     656332 kB\nInactive(file):   475072 kB\nUnevictable:        2844 kB\nMlocked:            2844 kB\nSwapTotal:             112345 kB\nSwapFree:             112340 kB\nDirty:                 0 kB\nWriteback:             0 kB\nAnonPages:        470460 kB\nMapped:           533244 kB\nShmem:               900 kB\nSlab:              81524 kB\nSReclaimable:      31544 kB\nSUnreclaim:        49980 kB\nKernelStack:       16016 kB\nPageTables:        24416 kB\nNFS_Unstable:          0 kB\nBounce:                0 kB\nWritebackTmp:          0 kB\nCommitLimit:      991448 kB\nCommitted_AS:   33535472 kB\nVmallocTotal:   258867136 kB\nVmallocUsed:           0 kB\nVmallocChunk:          0 kB\nCmaTotal:          40960 kB\nCmaFree:            2280 kB\nHugePages_Total:       0\nHugePages_Free:        0\nHugePages_Rsvd:        0\nHugePages_Surp:        0\nHugepagesize:       2048 kB"
	expected := &parser.Memory{
		MemTotal:  1982896,
		MemFree:   1258100,
		MemUsed:   724796,
		SwapTotal: 112345,
		SwapUsed:  5,
		SwapFree:  112340,
	}
	memory := parser.NewMemory()
	err := memory.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, memory)
}

func TestParseMemoryInvalidMemAvailable(t *testing.T) {
	t.Parallel()
	data := "MemTotal:        1982896 kB\nMemFree:          184732 kB\nMemAvailable:    0 kB\nBuffers:           30676 kB\nCached:          1103848 kB\nSwapCached:            0 kB\nActive:          1126396 kB\nInactive:         475672 kB\nActive(anon):     470064 kB\nInactive(anon):      600 kB\nActive(file):     656332 kB\nInactive(file):   475072 kB\nUnevictable:        2844 kB\nMlocked:            2844 kB\nSwapTotal:             0 kB\nSwapFree:              0 kB\nDirty:                 0 kB\nWriteback:             0 kB\nAnonPages:        470460 kB\nMapped:           533244 kB\nShmem:               900 kB\nSlab:              81524 kB\nSReclaimable:      31544 kB\nSUnreclaim:        49980 kB\nKernelStack:       16016 kB\nPageTables:        24416 kB\nNFS_Unstable:          0 kB\nBounce:                0 kB\nWritebackTmp:          0 kB\nCommitLimit:      991448 kB\nCommitted_AS:   33535472 kB\nVmallocTotal:   258867136 kB\nVmallocUsed:           0 kB\nVmallocChunk:          0 kB\nCmaTotal:          40960 kB\nCmaFree:            2280 kB\nHugePages_Total:       0\nHugePages_Free:        0\nHugePages_Rsvd:        0\nHugePages_Surp:        0\nHugepagesize:       2048 kB"
	expected := &parser.Memory{
		MemTotal: 1982896,
		MemFree:  184732,
		MemUsed:  1798164,
	}
	memory := parser.NewMemory()
	err := memory.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, memory)
}
