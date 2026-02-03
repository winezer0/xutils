package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func ToJSONLine(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func ToJSONBytesPretty(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// ToJSON 将任意 map 转换为格式化的 JSON 字符串（用于输出）
func ToJSON(v interface{}) string {
	data, _ := ToJSONBytesPretty(v)
	return string(data)
}

// LoadJSONBytes 从文件内容加载JSON数据
func LoadJSONBytes(content []byte, v interface{}) error {
	err := json.Unmarshal(content, v)
	return err
}

// LoadJSON 从文件加载JSON数据
func LoadJSON(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read json file error: %v", err)
	}
	if err := LoadJSONBytes(data, v); err != nil {
		return fmt.Errorf("parse json data error: %v", err)
	}
	return nil
}

// SaveJSON 保存JSON数据到文件
func SaveJSON(filePath string, v interface{}) error {
	data, err := ToJSONBytesPretty(v)
	if err != nil {
		return fmt.Errorf("serialization of JSON data failed: %v", err)
	}
	if err := SaveToFile(filePath, data); err != nil {
		return fmt.Errorf("failed to save the JSON file: %v", err)
	}
	return nil
}
