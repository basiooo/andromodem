package parser

import "strings"

type RootInfo struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}
type Root struct {
	IsRooted bool `json:"is_rooted"`
	*RootInfo
}

func ParseRoot(rawRoot string) Root {
	root := Root{}
	isRooted := !strings.Contains(rawRoot, "not found") && strings.TrimSpace(rawRoot) != ""
	root.IsRooted = isRooted
	if isRooted {
		rootInfo := strings.Split(strings.TrimSpace(rawRoot), ":")
		if len(rootInfo) == 2 {
			root.RootInfo = &RootInfo{
				Version: rootInfo[0],
				Name:    rootInfo[1],
			}
		}
	}
	return root
}
