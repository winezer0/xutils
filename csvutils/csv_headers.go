package csvutils

import (
	"encoding/csv"
	"fmt"
	"github.com/winezer0/xutils/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GetCSVHeaders 读取并解析 CSV 文件的首行作为表头信息。
// 该函数会自动处理分隔符检测（如果未指定）以及表头清洗（去重、去空格、处理 BOM）。
//
// 参数说明：
//   - filePath: CSV 文件的绝对或相对路径。
//   - delimiter: CSV 文件的分隔符。如果传入 0，函数将自动尝试检测分隔符。
//
// 返回值：
//   - []string: 清洗后的表头字符串切片。
//   - error: 执行过程中遇到的错误（如文件打开失败、读取失败）。
func GetCSVHeaders(filePath string, delimiter rune) ([]string, error) {
	// 处理空分隔符：默认设为逗号
	if delimiter == 0 { // 核心判断逻辑
		delimiter, _ = DetectCSVDelimiter(filePath)
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1 // 允许每行字段数不一致
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	headers = RepairHeaders(headers)

	return headers, nil
}

// GetCSVSHeadersMap 从多个 CSV 文件中收集每个文件的头部字段
// 返回：map[文件名][]表头、错误列表
// 不使用任何自定义类型，全部原生 Go 类型
func GetCSVSHeadersMap(csvFiles []string, delimiter rune, addPrefix bool) (map[string][]string, []error) {
	// 直接使用原生 map[string][]string
	fileHeaderMap := make(map[string][]string)
	var errors []error

	for _, filePath := range csvFiles {
		// 1. 获取纯文件名（例如从 "/path/to/data.csv" 获取 "data"）
		baseName := filepath.Base(filePath)
		fileName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

		// 2. 独立处理每个文件的分隔符
		currentDelimiter := delimiter
		if currentDelimiter == 0 {
			d, err := DetectCSVDelimiter(filePath)
			if err != nil {
				errors = append(errors, fmt.Errorf("detect delimiter failed for %s: %v", filePath, err))
				continue
			}
			currentDelimiter = d
		}

		// 3. 获取原始表头
		headers, err := GetCSVHeaders(filePath, currentDelimiter)
		if err != nil {
			errors = append(errors, fmt.Errorf("read headers failed for %s: %v", filePath, err))
			continue
		}

		// 4. 如果需要添加前缀，处理每一个表头
		var finalHeaders []string
		if addPrefix {
			finalHeaders = make([]string, 0, len(headers))
			for _, header := range headers {
				prefixed := fmt.Sprintf("%s.%s", fileName, header)
				finalHeaders = append(finalHeaders, prefixed)
			}
		} else {
			finalHeaders = headers
		}

		// 5. 存入原生 map：key=文件名，value=处理后的表头
		fileHeaderMap[fileName] = finalHeaders
	}

	return fileHeaderMap, errors
}

// GetCSVSHeaders 从多个 CSV 文件中收集所有唯一的头部字段。
// 【重构版：内部复用 GetCSVSHeadersMap，无重复代码】
func GetCSVSHeaders(csvFiles []string, delimiter rune, addPrefix bool) ([]string, []error) {
	// 1. 调用 map 版本获取所有文件的表头
	headerMap, errors := GetCSVSHeadersMap(csvFiles, delimiter, addPrefix)

	// 2. 遍历 map，收集并去重所有表头
	seenHeaders := make(map[string]bool)
	var allHeaders []string

	for _, headers := range headerMap {
		for _, header := range headers {
			if !seenHeaders[header] {
				seenHeaders[header] = true
				allHeaders = append(allHeaders, header)
			}
		}
	}

	return allHeaders, errors
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

// ShouldWriteHeader 判断是否需要向 CSV 文件写入表头
// 该函数用于在追加或写入 CSV 文件前，根据文件状态、覆盖模式以及新旧表头的对比结果，智能决策是否需要写入表头行。
// 决策逻辑：
// 空表头检查：如果传入的 header 为空，则无需写入。
// 覆盖模式：如果 overwrite 为 true，强制写入表头（通常意味着文件将被清空重写）。
// 文件存在性：如果文件不存在，必须写入表头以初始化文件。
// 表头一致性校验：
// 尝试读取现有文件的表头。
// 如果读取失败（如文件损坏或格式错误），为了安全起见，决定写入新表头。
// 如果现有表头与目标 header 不一致（顺序或内容不同），为了避免数据列错位，决定写入新表头（通常配合覆盖或修正逻辑）。
// 参数说明：
//
//	file: 目标 CSV 文件的路径。
//	header: 预期要写入的表头切片。
//	overwrite: 是否为覆盖模式。
//	delimiter: CSV 文件的分隔符（用于读取旧表头）。
//
// 返回值：
//
//	bool: true 表示需要写入表头，false 表示不需要（直接追加数据或跳过）。
func ShouldWriteHeader(file string, header []string, overwrite bool, delimiter rune) (should bool) {
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

	// csv头部长度不对
	if len(oldHeaders) != len(header) {
		fmt.Printf("file %s: header mismatch (old=%v, new=%v), will write new header\n", file, len(oldHeaders), len(header))
		return true
	}

	// csv头部顺序不对，写入，避免用户找不到有效头部开始行 清洗新表头，确保与旧表头使用相同的修复规则
	if !utils.SliceEqualStrict(RepairHeaders(oldHeaders), RepairHeaders(header)) {
		fmt.Printf("file %s: header mismatch (old=%v, new=%v), will write new header\n", file, oldHeaders, header)
		return true
	}

	// 其他情况下不需要写入头部
	return false
}

// FindIdenticalHeaders 传入 map[string][]string
// 返回：出现次数 >=2 的header（同一个文件内重复也累计次数）
func FindIdenticalHeaders(allHeaders map[string][]string) []string {
	// 统计所有 header 出现的总次数
	headerCount := make(map[string]int)

	// 遍历每个文件的表头
	for _, headers := range allHeaders {
		// 不去重！重复 header 直接累加计数
		for _, h := range headers {
			headerCount[h]++
		}
	}

	// 收集出现次数 >= 2 的 header
	var commonHeaders []string
	for header, count := range headerCount {
		if count >= 2 {
			commonHeaders = append(commonHeaders, header)
		}
	}

	return commonHeaders
}
