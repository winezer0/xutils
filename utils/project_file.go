package utils

import (
	"fmt"
	"github.com/winezer0/xutils/hashutils"
)

func GenProjectFileName(projectName, projectPath, toolName, suffix string) string {
	if len(suffix) == 0 {
		suffix = "cache"
	}
	return fmt.Sprintf("%s.%s.%s.%s", projectName, hashutils.GetStrHashShort(projectPath), toolName, suffix)
}
