package cmdutils

import (
	"fmt"
	"testing"
)

func TestParseDateTimeStrOutputs(t *testing.T) {
	cases := []string{
		"20251231:23:59:59",
		"2025-12-31:23:59:59",
		"2025/12/31:23:59:59",
		"20060102150405",
		"10h",
		"30m",
		"5s",
		"",
		"invalid",
	}
	fmt.Println("parse_date_time_str_outputs:")
	for _, c := range cases {
		ts, err := ParseExitTime(c)
		if err != nil {
			fmt.Println(c, "=>", "error:", err)
		} else if ts == nil {
			fmt.Println(c, "=>", "nil")
		} else {
			fmt.Println(c, "=>", ts.String())
		}
	}
}
