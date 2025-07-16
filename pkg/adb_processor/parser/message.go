package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	regRow      = regexp.MustCompile(`Row:([\s\S]*?) address`)
	regAddress  = regexp.MustCompile(`address=([^,]+)`)
	regBody     = regexp.MustCompile(`body=([\s\S]*?), date`)
	regDate     = regexp.MustCompile(`date=(\d+)`)
	wibLocation *time.Location
)

func init() {
	var err error
	wibLocation, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// fallback ke UTC kalau gagal load lokasi
		wibLocation = time.UTC
	}
}

type Message struct {
	Row     int    `json:"row"`
	Address string `json:"address"`
	Body    string `json:"body"`
	Date    string `json:"date"`
}

func NewMessage() IParser {
	return &Message{}
}

func (m *Message) IsEmpty() bool {
	return m.Address == "" && m.Body == "" && m.Date == ""
}

func (m *Message) formatDate(timestamp string) string {
	ms, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ""
	}
	// Unix expects seconds, timestamp input in milliseconds
	t := time.Unix(ms/1000, 0).In(wibLocation)
	return t.Format("2006-01-02 15:04:05")
}

func (m *Message) Parse(rawData string) error {
	m.Row = extractInt(regRow, rawData)
	m.Address = extractString(regAddress, rawData)
	m.Body = extractString(regBody, rawData)

	rawDate := extractString(regDate, rawData)
	if rawDate != "" {
		m.Date = m.formatDate(rawDate)
	}
	return nil
}

func extractString(r *regexp.Regexp, data string) string {
	match := r.FindStringSubmatch(data)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func extractInt(r *regexp.Regexp, data string) int {
	strVal := extractString(r, data)
	val, err := strconv.Atoi(strVal)
	if err != nil {
		return 0
	}
	return val
}
