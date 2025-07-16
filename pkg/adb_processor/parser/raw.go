package parser

import "strings"

type RawParser struct {
	Result string
}

func NewRawParser() IParser {
	return &RawParser{}
}

func (d *RawParser) Parse(data string) error {
	d.Result = strings.TrimSpace(data)
	return nil
}
