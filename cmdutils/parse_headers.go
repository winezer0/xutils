package cmdutils

import (
	"fmt"
	"strings"
)

// ParseHeaders 从形如 key=value 的切片解析请求头(中文注释)
func ParseHeaders(headers []string) map[string]string {
	headerMaps := make(map[string]string)
	for _, item := range headers {
		if item == "" {
			continue
		}
		kv := strings.SplitN(item, "=", 2)
		if len(kv) != 2 {
			fmt.Printf("Invalid header: %s\n", item)
			continue
		}
		headerMaps[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return headerMaps
}
