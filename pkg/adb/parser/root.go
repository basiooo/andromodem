package parser

import "strings"

type RootDetail struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

type Root struct {
	IsRooted bool `json:"is_rooted"`
	*RootDetail
}

func NewRoot(rawRoot string) *Root {
	root := &Root{}
	root.IsRooted = root.isRooted(rawRoot)
	if root.IsRooted {
		root.RootDetail = root.extractRootDetail(rawRoot)
	}
	return root
}

func (r *Root) isRooted(rawRoot string) bool {
	return !strings.Contains(rawRoot, "not found") && strings.TrimSpace(rawRoot) != ""
}

func (r *Root) extractRootDetail(rawRoot string) *RootDetail {
	rootDetail := strings.Split(strings.TrimSpace(rawRoot), ":")
	if len(rootDetail) == 2 {
		return &RootDetail{
			Version: rootDetail[0],
			Name:    rootDetail[1],
		}
	}
	return nil
}
