package parser_test

import (
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSignalStrength(t *testing.T) {
	t.Parallel()
	data := `{mCdma=CellSignalStrengthCdma: cdmaDbm=2147483647 cdmaEcio=2147483647 evdoDbm=2147483647 evdoEcio=2147483647 evdoSnr=2147483647 level=0,mGsm=CellSignalStrengthGsm: rssi=2147483647 ber=2147483647 mTa=2147483647 mLevel=0,mWcdma=CellSignalStrengthWcdma: ss=2147483647 ber=2147483647 rscp=2147483647 ecno=2147483647 level=0,mTdscdma=CellSignalStrengthTdscdma: rssi=2147483647 ber=2147483647 rscp=2147483647 level=0,mLte=CellSignalStrengthLte: rssi=-65 rsrp=-95 rsrq=-11 rssnr=8 cqiTableIndex=2147483647 cqi=2147483647 ta=2147483647 level=4 parametersUseForLevel=0,mNr=CellSignalStrengthNr:{ csiRsrp = 2147483647 csiRsrq = 2147483647 csiCqiTableIndex = 2147483647 csiCqiReport = [] ssRsrp = 2147483647 ssRsrq = 2147483647 ssSinr = 2147483647 level = 0 parametersUseForLevel = 0 timingAdvance = 2147483647 },primary=CellSignalStrengthLte}`
	expected := &parser.SignalStrength{
		CellSignalStrengthLte: &parser.CellSignalStrengthLte{
			Rssi:                  "-65",
			Rsrp:                  "-95",
			Rsrq:                  "-11",
			Rssnr:                 "8",
			CqiTableIndex:         "2147483647",
			Cqi:                   "2147483647",
			Ta:                    "2147483647",
			Level:                 "4",
			ParametersUseForLevel: "0",
		},
	}
	signalStrength := parser.NewSignalStrength()
	err := signalStrength.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, signalStrength)
}

func TestParseSignalStrengOlderAndroidVersion(t *testing.T) {
	t.Parallel()
	data := "lte"
	expected := &parser.SignalStrength{
		CellSignalStrengthLte: &parser.CellSignalStrengthLte{},
	}
	signalStrength := parser.NewSignalStrength()
	err := signalStrength.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, signalStrength)
}

func TestParseSignalStrengthNoSim(t *testing.T) {
	t.Parallel()
	data := ""
	expected := &parser.SignalStrength{}
	signalStrength := parser.NewSignalStrength()
	err := signalStrength.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, signalStrength)
}
