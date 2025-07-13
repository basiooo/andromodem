package parser_test

import (
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSmsInbox(t *testing.T) {
	t.Parallel()
	data := `Row: 433 address=XL-Axiata, body=(XL) Kuota utama Pkt Xtra Combo Mini Anda sdh habis.Saat ini berlaku tarif dasar internet. Aktifkan lagi paketnya di aplikasi MyXL., date=1712904715243
	Row: 440 address=3TopUp, body=Isi pulsa langsung dapat kuota 24 jam
	-5rb,500MB 1hr
	-10rb,1GB 1hr
	-20rb,1.5GB 2hr
	-50rb,3GB 3hr
	-100rb,5GB 3hr
	Di toko/app favoritmu atau bit.ly/1sipls4, date=1712990828349
	`
	expected := &parser.Inbox{
		Messages: []parser.Message{
			{
				Row:     433,
				Address: "XL-Axiata",
				Body:    "(XL) Kuota utama Pkt Xtra Combo Mini Anda sdh habis.Saat ini berlaku tarif dasar internet. Aktifkan lagi paketnya di aplikasi MyXL.",
				Date:    "2024-04-12 13:51:55",
			},
			{
				Row:     440,
				Address: "3TopUp",
				Body:    "Isi pulsa langsung dapat kuota 24 jam\n\t-5rb,500MB 1hr\n\t-10rb,1GB 1hr\n\t-20rb,1.5GB 2hr\n\t-50rb,3GB 3hr\n\t-100rb,5GB 3hr\n\tDi toko/app favoritmu atau bit.ly/1sipls4",
				Date:    "2024-04-13 13:47:08",
			},
		},
	}
	inbox := parser.NewInbox()
	err := inbox.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, inbox)
}
func TestParseSmsInboxEmpty(t *testing.T) {
	t.Parallel()
	data := `no result found.`
	expected := &parser.Inbox{}
	inbox := parser.NewInbox()
	err := inbox.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, inbox)
}

func BenchmarkParseSmsInbox(b *testing.B) {
	data := `Row: 433 address=XL-Axiata, body=(XL) Kuota utama Pkt Xtra Combo Mini Anda sdh habis.Saat ini berlaku tarif dasar internet. Aktifkan lagi paketnya di aplikasi MyXL., date=1712904715243
	Row: 440 address=3TopUp, body=Isi pulsa langsung dapat kuota 24 jam
	-5rb,500MB 1hr
	-10rb,1GB 1hr
	-20rb,1.5GB 2hr
	-50rb,3GB 3hr
	-100rb,5GB 3hr
	Di toko/app favoritmu atau bit.ly/1sipls4, date=1712990828349
	`
	for i := 0; i < b.N; i++ {
		inbox := parser.NewInbox()
		_ = inbox.Parse(data)
	}
}

func BenchmarkParseSmsInboxEmpty(b *testing.B) {
	data := `no result found.`
	for i := 0; i < b.N; i++ {
		inbox := parser.NewInbox()
		_ = inbox.Parse(data)
	}
}
