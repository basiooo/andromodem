package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseSmsInbox(t *testing.T) {
	data := `Row: 433 address=XL-Axiata, body=(XL) Kuota utama Pkt Xtra Combo Mini Anda sdh habis.Saat ini berlaku tarif dasar internet. Aktifkan lagi paketnya di aplikasi MyXL., date=1712904715243
	Row: 440 address=3TopUp, body=Isi pulsa langsung dapat kuota 24 jam
	-5rb,500MB 1hr
	-10rb,1GB 1hr
	-20rb,1.5GB 2hr
	-50rb,3GB 3hr
	-100rb,5GB 3hr
	Di toko/app favoritmu atau bit.ly/1sipls4, date=1712990828349
	`
	expected := []parser.SMSInbox{{
		Row:     "433",
		Address: "XL-Axiata",
		Body:    "(XL) Kuota utama Pkt Xtra Combo Mini Anda sdh habis.Saat ini berlaku tarif dasar internet. Aktifkan lagi paketnya di aplikasi MyXL.",
		Date:    "2024-04-12 13:51:55",
	}, {
		Row:     "440",
		Address: "3TopUp",
		Body:    "Isi pulsa langsung dapat kuota 24 jam\n\t-5rb,500MB 1hr\n\t-10rb,1GB 1hr\n\t-20rb,1.5GB 2hr\n\t-50rb,3GB 3hr\n\t-100rb,5GB 3hr\n\tDi toko/app favoritmu atau bit.ly/1sipls4",
		Date:    "2024-04-13 13:47:08",
	}}
	smsInboxs, err := parser.NewSMSInbox(data)
	actual := *smsInboxs
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSmsInboxNotPermission(t *testing.T) {
	data := `Error while accessing provider:sms
	java.lang.SecurityException: Permission Denial	
	`
	_, err := parser.NewSMSInbox(data)
	assert.NotNil(t, err)
}
