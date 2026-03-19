package csvutils

import (
	"fmt"
	"github.com/winezer0/xutils/utils"
	"sort"
)

// 检查是否需要写入表头
func shouldWriteHeader(file string, header []string, overwrite bool, delimiter rune) (should bool) {
	// 没有头部自然不需要写入头部
	if len(header) == 0 {
		return false
	}

	// 覆盖模式需要写入csv头部
	if overwrite {
		return true
	}

	// 文件不存在 肯定要写入头部
	if !utils.FileExists(file) {
		return true
	}

	// 检查其他模式
	oldHeaders, err := GetCSVHeaders(file, delimiter)
	// 读取csv头失败，写入，避免用户找不到有效头部开始行
	if err != nil {
		fmt.Printf("file %s: read old header failed, err=%v, will write new header\n", file, err)
		return true
	}

	// csv头部顺序不对，写入，避免用户找不到有效头部开始行
	if !utils.SliceEqualStrict(oldHeaders, header) {
		fmt.Printf("file %s: header mismatch (old=%v, new=%v), will write new header\n", file, oldHeaders, header)
		return true
	}

	// 其他情况下不需要写入头部
	return false
}

// GetCSVHeaderFromDicts 从字典列表中获取待写入的 CSV 表头
// 核心逻辑：
// 1. 若传入的 header 非空，直接返回（优先使用业务指定的表头）
// 2. 若 header 为空，根据 deep 模式生成表头：
//   - deep=false：仅从第一个字典提取键（效率高，适用于字典结构统一的场景）
//   - deep=true：遍历所有字典，收集所有出现过的键（适用于字典结构不统一的场景）
//
// 3. 自动生成的表头会排序，保证顺序稳定
// 参数：
//
//	dicts - 字典列表（每个字典代表一行数据）
//	header   - 业务指定的表头（nil/空切片则自动生成）
//	deep     - 是否深度遍历所有字典收集键（false=仅第一个字典，true=所有字典）
//
// 返回值：
//
//	最终用于写入 CSV 的表头切片
//	错误（仅当 dicts 为空时返回）
func GetCSVHeaderFromDicts(dicts []map[string]interface{}, header []string, deep bool) ([]string, error) {
	// 1. 优先使用业务指定的表头
	if len(header) > 0 {
		// 拷贝切片，避免外部修改影响内部逻辑
		usedHeader := make([]string, len(header))
		copy(usedHeader, header)
		return usedHeader, nil
	}

	// 2. 边界处理：字典列表为空时无法生成表头
	if len(dicts) == 0 {
		return nil, fmt.Errorf("dicts is empty, cannot generate header")
	}

	// 3. 初始化键集合（用 map 去重）
	keySet := make(map[string]struct{})

	// 4. 根据 deep 模式收集键
	if deep {
		// deep=true：遍历所有字典，收集所有键
		for _, dict := range dicts {
			for k := range dict {
				keySet[k] = struct{}{}
			}
		}
	} else {
		// deep=false：仅从第一个字典收集键（默认模式，效率高）
		firstDict := dicts[0]
		for k := range firstDict {
			keySet[k] = struct{}{}
		}
	}

	// 5. 将去重后的键转为切片并排序（保证顺序稳定）
	usedHeader := make([]string, 0, len(keySet))
	for k := range keySet {
		usedHeader = append(usedHeader, k)
	}
	sort.Strings(usedHeader)

	return usedHeader, nil
}
