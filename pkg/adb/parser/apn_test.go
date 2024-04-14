package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseApn(t *testing.T) {
	data := `Row: 0 _id=2943, name=XL Unlimited, numeric=51011, mcc=510, mnc=11, carrier_id=-1, apn=xlunlimited, user=, server=, password=, proxy=202.152.240.50, port=8080, mmsproxy=, mmsport=, mmsc=, authtype=-1, type=default,supl, current=1, protocol=IP, roaming_protocol=IP, carrier_enabled=1, bearer=0, bearer_bitmask=0, network_type_bitmask=0, lingering_network_type_bitmask=0, mvno_type=, mvno_match_data=, sub_id=-1, profile_id=0, modem_cognitive=0, max_conns=0, wait_time=0, max_conns_time=0, mtu=0, mtu_v4=0, mtu_v6=0, edited=0, user_visible=1, user_editable=1, owned_by=1, apn_set_id=0, skip_464xlat=-1, always_on=0`
	expected := parser.Apn{
		Name:    "XL Unlimited",
		ApnName: "xlunlimited",
	}
	apn := parser.NewApn(data)
	actual := *apn
	assert.Equal(t, expected, actual)
}

func TestParseApnNoPermission(t *testing.T) {
	data := `Error while accessing provider:telephony
	java.lang.SecurityException: No permission to write APN settings
		at android.os.Parcel.createException(Parcel.java:2071)
		at android.os.Parcel.readException(Parcel.java:2039)
		at android.database.DatabaseUtils.readExceptionFromParcel(DatabaseUtils.java:188)
		at android.database.DatabaseUtils.readExceptionFromParcel(DatabaseUtils.java:140)
		at android.content.ContentProviderProxy.query(ContentProviderNative.java:423)
		at com.android.commands.content.Content$QueryCommand.onExecute(Content.java:619)
		at com.android.commands.content.Content$Command.execute(Content.java:469)
		at com.android.commands.content.Content.main(Content.java:690)
		at com.android.internal.os.RuntimeInit.nativeFinishInit(Native Method)
		at com.android.internal.os.RuntimeInit.main(RuntimeInit.java:342)`
	expected := parser.Apn{
		Name:    "",
		ApnName: "",
	}
	apn := parser.NewApn(data)
	actual := *apn
	assert.Equal(t, expected, actual)
}
