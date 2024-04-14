package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseCarriers(t *testing.T) {
	rawData := parser.RawCarriers{
		RawConnectionsState: "-1\n2",
		RawCarriersName:     "3-SinyalKuatLuas,XL Axiata",
		RawSignalsStrength: `{mCdma=CellSignalStrengthCdma: cdmaDbm=2147483647 cdmaEcio=2147483647 evdoDbm=2147483647 evdoEcio=2147483647 evdoSnr=2147483647 level=0,mGsm=CellSignalStrengthGsm: rssi=2147483647 ber=2147483647 mTa=2147483647 mLevel=0,mWcdma=CellSignalStrengthWcdma: ss=2147483647 ber=2147483647 rscp=2147483647 ecno=2147483647 level=0,mTdscdma=CellSignalStrengthTdscdma: rssi=2147483647 ber=2147483647 rscp=2147483647 level=0,mLte=CellSignalStrengthLte: rssi=-65 rsrp=-95 rsrq=-11 rssnr=8 cqiTableIndex=2147483647 cqi=2147483647 ta=2147483647 level=4 parametersUseForLevel=0,mNr=CellSignalStrengthNr:{ csiRsrp = 2147483647 csiRsrq = 2147483647 csiCqiTableIndex = 2147483647 csiCqiReport = [] ssRsrp = 2147483647 ssRsrq = 2147483647 ssSinr = 2147483647 level = 0 parametersUseForLevel = 0 timingAdvance = 2147483647 },primary=CellSignalStrengthLte}
		{mCdma=CellSignalStrengthCdma: cdmaDbm=2147483647 cdmaEcio=2147483647 evdoDbm=2147483647 evdoEcio=2147483647 evdoSnr=2147483647 level=0,mGsm=CellSignalStrengthGsm: rssi=2147483647 ber=2147483647 mTa=2147483647 mLevel=0,mWcdma=CellSignalStrengthWcdma: ss=2147483647 ber=2147483647 rscp=2147483647 ecno=2147483647 level=0,mTdscdma=CellSignalStrengthTdscdma: rssi=2147483647 ber=2147483647 rscp=2147483647 level=0,mLte=CellSignalStrengthLte: rssi=-61 rsrp=-91 rsrq=-10 rssnr=10 cqiTableIndex=2147483647 cqi=2147483647 ta=2147483647 level=4 parametersUseForLevel=0,mNr=CellSignalStrengthNr:{ csiRsrp = 2147483647 csiRsrq = 2147483647 csiCqiTableIndex = 2147483647 csiCqiReport = [] ssRsrp = 2147483647 ssRsrq = 2147483647 ssSinr = 2147483647 level = 0 parametersUseForLevel = 0 timingAdvance = 2147483647 },primary=CellSignalStrengthLte}`,
	}
	expected := []parser.Carrier{
		{
			Name:            "3-SinyalKuatLuas",
			ConnectionState: parser.DataUnknown.String(),
			SimSlot:         1,
			SignalStrength: parser.SignalStrength{
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
			},
		},
		{
			Name:            "XL Axiata",
			ConnectionState: parser.DataConnected.String(),
			SimSlot:         2,
			SignalStrength: parser.SignalStrength{
				CellSignalStrengthLte: &parser.CellSignalStrengthLte{
					Rssi:                  "-61",
					Rsrp:                  "-91",
					Rsrq:                  "-10",
					Rssnr:                 "10",
					CqiTableIndex:         "2147483647",
					Cqi:                   "2147483647",
					Ta:                    "2147483647",
					Level:                 "4",
					ParametersUseForLevel: "0",
				},
			},
		},
	}
	carriersInfo := parser.NewCarriers(rawData)
	actual := *carriersInfo
	assert.Equal(t, expected, actual)
}
