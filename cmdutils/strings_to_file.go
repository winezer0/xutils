package cmdutils

import (
	"fmt"
	"github.com/winezer0/xutils/hashutils"
	"github.com/winezer0/xutils/utils"
	"strings"
)

// ConvertStringsToFiles 智能处理输入列表，将其转换为有效的文件路径列表。
//
// 行为逻辑：
//  1. 预处理：若 clean 为 true，则对 inputs 进行去重和去除空字符串处理（保留大小写敏感）。
//  2. 检测：检查列表中哪些是已存在的文件路径，哪些是纯文本字符串。
//  3. 分支处理：
//     - 若全为有效文件：直接返回原路径列表。
//     - 若包含纯文本：将所有纯文本内容排序后写入一个临时文件，并将该临时文件路径
//     与原有的有效文件路径合并返回。
//
// 参数:
//   - inputs: 输入列表，元素可能是现有的文件路径，也可能是待处理的纯文本字符串。
//   - clean: 是否执行清理操作（去重、去空）。注意：此操作不忽略大小写。
//
// 返回:
//   - files: 最终有效的文件路径列表。如果生成了新文件，列表中将包含该新文件的路径。
//   - newly: 如果生成了新文件，返回其完整路径；否则返回空字符串。
func ConvertStringsToFiles(inputs []string, clean bool) (files []string, newly string) {
	// ... (函数实现保持不变)
	if len(inputs) == 0 {
		return inputs, ""
	}

	// 进行去重处理
	if clean {
		inputs = utils.UniqueSlice(inputs, false, true)
	}

	newFile := ""
	validFilePaths, plainStrings := utils.CheckFilesExist(inputs)
	if len(plainStrings) > 0 {
		plainStrings = utils.SortSlice(plainStrings)
		tempFile := fmt.Sprintf("%s.temp", hashutils.GetStrHash(strings.Join(plainStrings, "||")))
		if err := utils.WriteLines(tempFile, plainStrings, true); err == nil {
			validFilePaths = append(validFilePaths, tempFile)
			newFile = tempFile
		}
	}

	// 返回仅包含新生成文件路径的切片
	return validFilePaths, newFile
}
