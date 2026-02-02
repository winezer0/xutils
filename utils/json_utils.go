package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadJSON 从文件加载JSON数据
func LoadJSON(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read json file error: %v", err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("parse json data error: %v", err)
	}

	return nil
}

// SaveJSON 保存JSON数据到文件
func SaveJSON(filePath string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("serialization of JSON data failed: %v", err)
	}
	if err := SaveToFile(filePath, data); err != nil {
		return fmt.Errorf("failed to save the JSON file: %v", err)
	}
	return nil
}

// ToJson 将任意 map 转换为格式化的 JSON 字符串（用于输出）
func ToJson(v interface{}) string {
	return string(ToJsonBytes(v))
}

// ToJsonBytes  将任意 map 转换为格式化的 JSON 字符串（用于输出）
func ToJsonBytes(v interface{}) []byte {
	data, _ := json.MarshalIndent(v, "", "  ")
	return data
}

// ToJsonLine 将任意值转换为单行 JSON 字符串
func ToJsonLine(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
