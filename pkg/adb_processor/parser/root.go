package parser

import (
	"strings"
)

type RootDetail struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

type Root struct {
	IsRooted                  bool `json:"is_rooted"`
	SuperUserAllowShellAccess bool `json:"super_user_allow_shell_access"`
	RootDetail
}

func NewRoot() IParser {
	return &Root{}
}

func (r *Root) Parse(rawData string) error {
	trimmed := strings.TrimSpace(rawData)
	r.IsRooted = trimmed != "" && !strings.Contains(trimmed, "not found")

	if r.IsRooted {
		r.RootDetail = parseRootDetail(trimmed)
	}

	return nil
}

func parseRootDetail(rawData string) RootDetail {
	rootDetail := RootDetail{}
	rawData = strings.TrimSpace(rawData)
	detail := strings.Split(rawData, ":")
	if len(detail) != 2 {
		detail = strings.Split(rawData, " ")
	}
	if len(detail) == 2 {
		rootDetail.Version = detail[0]
		rootDetail.Name = detail[1]
	} else {
		rootDetail.Name = rawData
	}
	return rootDetail
}
