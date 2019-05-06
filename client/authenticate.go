package client

import (
	"fmt"
	"regexp"
	"strings"
)

func parseAuthenticate(txt string) (map[string]string, error) {
	// https://tools.ietf.org/html/rfc2617#page-8
	values := make(map[string]string)
	re := regexp.MustCompile(`(\w+)\s*=\s*"(.*)"`)
	for _, blob := range strings.Split(txt, ",") {
		m := re.FindStringSubmatch(blob)
		if len(m) != 3 {
			return nil, fmt.Errorf("Invalid header: %v", blob)
		}
		values[m[1]] = m[2]
	}
	return values, nil
}
