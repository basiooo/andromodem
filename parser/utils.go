package parser

import (
	"fmt"
	"regexp"
)

func GetPropValueByName(prop, name string) string {
	regexPattern := fmt.Sprintf(`\[%s\]: \[(.*?)\]`, name)
	r := regexp.MustCompile(regexPattern)

	match := r.FindStringSubmatch(prop)

	if len(match) < 2 {
		return ""
	}

	return match[1]
}
