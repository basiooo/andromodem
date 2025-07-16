package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

type Storage struct {
	DataTotal   int `json:"data_total,omitempty"`
	DataFree    int `json:"data_free,omitempty"`
	DataUsed    int `json:"data_used,omitempty"`
	SystemTotal int `json:"system_total,omitempty"`
	SystemFree  int `json:"system_free,omitempty"`
	SystemUsed  int `json:"system_used,omitempty"`
}

var re = regexp.MustCompile(`(\d+)K\s*/\s*(\d+)K`)

func NewStorage() IParser {
	return &Storage{}
}

func parseStorageLine(storageData map[string]string, label string) (free, total int, err error) {
	value := getParseListDataValue(storageData, label)
	matches := re.FindStringSubmatch(value)
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("invalid format for %s: %q", label, value)
	}

	free, err = strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse free space in %s: %w", label, err)
	}

	total, err = strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse total space in %s: %w", label, err)
	}

	return free, total, nil
}
func (s *Storage) Parse(rawData string) error {
	storageData := parseListData(rawData, []string{
		"Data-Free",
		"System-Free",
	}, false)

	if free, total, err := parseStorageLine(storageData, "Data-Free"); err == nil {
		s.DataFree = free
		s.DataTotal = total
		s.DataUsed = total - free
	}

	if free, total, err := parseStorageLine(storageData, "System-Free"); err == nil {
		s.SystemFree = free
		s.SystemTotal = total
		s.SystemUsed = total - free
	}
	return nil
}
