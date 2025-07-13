package parser

import (
	"strings"
)

type BusyboxCheck struct {
	Exist bool
}

func NewBusyboxCheck() IParser {
	return &BusyboxCheck{}
}

func (b *BusyboxCheck) Parse(rawData string) error {
	b.Exist = strings.Contains(rawData, "busybox")
	return nil
}
