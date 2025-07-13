package parser

import (
	"strconv"
)

type Memory struct {
	MemTotal  int `json:"mem_total"`
	MemFree   int `json:"mem_free"`
	MemUsed   int `json:"mem_used"`
	SwapTotal int `json:"swap_total"`
	SwapFree  int `json:"swap_free"`
	SwapUsed  int `json:"swap_used"`
}

func NewMemory() IParser {
	return &Memory{}
}

func (m *Memory) Parse(rawData string) error {
	fields := []string{
		"MemTotal",
		"MemAvailable",
		"MemFree",
		"SwapTotal",
		"SwapFree",
	}
	memoryData := parseListData(rawData, fields, true)

	memTotal := getParseListDataValue(memoryData, "MemTotal")
	memFree := getParseListDataValue(memoryData, "MemAvailable")
	swapTotal := getParseListDataValue(memoryData, "SwapTotal")
	swapFree := getParseListDataValue(memoryData, "SwapFree")

	if val, _ := strconv.Atoi(memFree); val == 0 {
		memFree = getParseListDataValue(memoryData, "MemFree")
	}

	var err error
	if m.MemTotal, err = strconv.Atoi(memTotal); err != nil {
		return err
	}
	if m.MemTotal, err = strconv.Atoi(memTotal); err != nil {
		return err
	}
	if m.MemFree, err = strconv.Atoi(memFree); err != nil {
		return err
	}
	if m.SwapTotal, err = strconv.Atoi(swapTotal); err != nil {
		return err
	}
	if m.SwapFree, err = strconv.Atoi(swapFree); err != nil {
		return err
	}

	m.MemUsed = m.MemTotal - m.MemFree
	m.SwapUsed = m.SwapTotal - m.SwapFree
	return nil
}
