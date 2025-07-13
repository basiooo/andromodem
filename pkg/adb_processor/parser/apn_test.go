package parser_test

import (
	"encoding/json"
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseApn(t *testing.T) {
	t.Parallel()
	data := `Row: 0 _id=2943, name=XL Unlimited, numeric=51011, mcc=510, mnc=11, carrier_id=-1, apn=xlunlimited, user=, server=, password=, proxy=202.152.240.50, port=8080, mmsproxy=, mmsport=, mmsc=, authtype=-1, type=default,supl, current=1, protocol=IP, roaming_protocol=IP, carrier_enabled=1, bearer=0, bearer_bitmask=0, network_type_bitmask=0, lingering_network_type_bitmask=0, mvno_type=, mvno_match_data=, sub_id=-1, profile_id=0, modem_cognitive=0, max_conns=0, wait_time=0, max_conns_time=0, mtu=0, mtu_v4=0, mtu_v6=0, edited=0, user_visible=1, user_editable=1, owned_by=1, apn_set_id=0, skip_464xlat=-1, always_on=0`
	expected := &parser.Apn{
		Name:     "XL Unlimited",
		ApnName:  "xlunlimited",
		Proxy:    "202.152.240.50",
		Port:     "8080",
		Username: "",
		Password: "",
		Server:   "",
		MCC:      "510",
		MNC:      "11",
		Type:     "default,supl",
		Protocol: "IP",
	}
	apn := parser.NewApn()
	err := apn.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, apn)
}

func TestApn_MarshalJSON(t *testing.T) {
	t.Parallel()
	data := `Row: 0 _id=2943, name=XL Unlimited, numeric=51011, mcc=510, mnc=11, carrier_id=-1, apn=xlunlimited, user=, server=, password=, proxy=202.152.240.50, port=8080, mmsproxy=, mmsport=, mmsc=, authtype=-1, type=default,supl, current=1, protocol=IP, roaming_protocol=IP, carrier_enabled=1, bearer=0, bearer_bitmask=0, network_type_bitmask=0, lingering_network_type_bitmask=0, mvno_type=, mvno_match_data=, sub_id=-1, profile_id=0, modem_cognitive=0, max_conns=0, wait_time=0, max_conns_time=0, mtu=0, mtu_v4=0, mtu_v6=0, edited=0, user_visible=1, user_editable=1, owned_by=1, apn_set_id=0, skip_464xlat=-1, always_on=0`
	apn := parser.NewApn()
	err := apn.Parse(data)
	assert.NoError(t, err)
	result, err := json.Marshal(apn)
	assert.NoError(t, err)
	expected := `{"name":"XL Unlimited","apn":"xlunlimited","proxy":"202.152.240.50","port":"8080","mms_proxy":"","mms_port":"","username":"","password":"","server":"","mcc":"510","mnc":"11","type":"default,supl","protocol":"IP"}`
	assert.Equal(t, expected, string(result))
}

func BenchmarkParseApn(b *testing.B) {
	data := `Row: 0 _id=2943, name=XL Unlimited, numeric=51011, mcc=510, mnc=11, carrier_id=-1, apn=xlunlimited, user=, server=, password=, proxy=202.152.240.50, port=8080, mmsproxy=, mmsport=, mmsc=, authtype=-1, type=default,supl, current=1, protocol=IP, roaming_protocol=IP, carrier_enabled=1, bearer=0, bearer_bitmask=0, network_type_bitmask=0, lingering_network_type_bitmask=0, mvno_type=, mvno_match_data=, sub_id=-1, profile_id=0, modem_cognitive=0, max_conns=0, wait_time=0, max_conns_time=0, mtu=0, mtu_v4=0, mtu_v6=0, edited=0, user_visible=1, user_editable=1, owned_by=1, apn_set_id=0, skip_464xlat=-1, always_on=0`
	for i := 0; i < b.N; i++ {
		apn := parser.NewApn()
		_ = apn.Parse(data)
	}
}
