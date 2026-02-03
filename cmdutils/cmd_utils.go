package cmdutils

import "strings"

func ContainsComma(s string) bool { return strings.Contains(s, ",") }

func SplitComma(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
