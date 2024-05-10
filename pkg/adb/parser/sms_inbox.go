package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type SMSInbox struct {
	Row     string `json:"row"`
	Address string `json:"address"`
	Body    string `json:"body"`
	Date    string `json:"date"`
}

func NewSMSInbox(rawSmsInbox string) (*[]SMSInbox, error) {
	splitSms := splitSmsinbox(rawSmsInbox)
	smsInboxs := parseSmsInbox(splitSms)
	if len(*smsInboxs) == 0 && !isEmptySMS(rawSmsInbox) {
		if isErrorReqRootPermission(rawSmsInbox) {
			return nil, errors.New("cannot get sms list. please allow root permission")
		}

		if isErrorReqRootDevice(rawSmsInbox) {
			return nil, errors.New("cannot get sms list without root")
		}
	}
	return smsInboxs, nil
}

func isErrorReqRootDevice(rawSmsInbox string) bool {
	return strings.Contains(strings.ToLower(rawSmsInbox), "error while accessing provider") || !strings.Contains(rawSmsInbox, "Row: ")
}

func isEmptySMS(rawSmsInbox string) bool {
	return strings.TrimSpace(strings.ToLower(rawSmsInbox)) == "no result found."
}

func isErrorReqRootPermission(rawSmsInbox string) bool {
	return strings.Contains(strings.ToLower(rawSmsInbox), "permission")
}

func splitSmsinbox(rawSmsInbox string) []string {
	raw := strings.Split(rawSmsInbox, "Row:")
	finalRaw := make([]string, 0, len(raw)-1)
	for _, s := range raw {
		if s == "" {
			continue
		}
		finalRaw = append(finalRaw, "Row:"+s)
	}
	return finalRaw
}

func formatDate(timestamp string) string {
	timestampMs, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		logrus.WithField("location", "sms_inbox.GetSmsInbox").Errorf("formatDate(): failed format date: %v", err)
		return ""
	}
	t := time.Unix(timestampMs/1000, 0).UTC()
	jakartaLocation, _ := time.LoadLocation("Asia/Jakarta")
	localTime := t.In(jakartaLocation)
	return localTime.Format("2006-01-02 15:04:05")
}

func parseSmsInbox(rawSmsInboxs []string) *[]SMSInbox {
	reRow := regexp.MustCompile(`Row:([\s\S]*?) address`)
	reAddress := regexp.MustCompile(`address=([^,]+)`)
	reBody := regexp.MustCompile(`body=([\s\S]*?), date`)
	reDate := regexp.MustCompile(`date=(\d+)`)
	smsInboxs := make([]SMSInbox, 0, len(rawSmsInboxs))
	for _, input := range rawSmsInboxs {
		smsInbox := SMSInbox{}
		row := reRow.FindStringSubmatch(input)
		address := reAddress.FindStringSubmatch(input)
		body := reBody.FindStringSubmatch(input)
		rawDate := reDate.FindStringSubmatch(input)
		if len(address) == 2 {
			smsInbox.Address = address[1]
		}
		if len(row) == 2 {
			smsInbox.Row = strings.TrimSpace(row[1])
		}
		if len(body) == 2 {
			smsInbox.Body = body[1]
		}
		if len(rawDate) == 2 {
			smsInbox.Date = formatDate(rawDate[1])
			smsInboxs = append(smsInboxs, smsInbox)
		}
	}
	return &smsInboxs
}
