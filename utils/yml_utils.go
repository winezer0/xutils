package utils

import (
	"fmt"
	"github.com/winezer0/xutils/logging"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadYAMLBytes 从文件内容中加载YAML数据 v 一般是地址对象
func LoadYAMLBytes(content []byte, v interface{}) error {
	if err := yaml.Unmarshal(content, v); err != nil {
		return err
	}
	return nil
}

// LoadYAML 从文件加载YAML数据
func LoadYAML(filePath string, v interface{}) error {
	logging.Debugf("Loading YAML file...: %s", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read the YAML file: %v", err)
	}
	logging.Debugf("Size of the YAML file: %d byte", len(data))
	if err := LoadYAMLBytes(data, v); err != nil {
		return fmt.Errorf("failed to parse YAML data: %v", err)
	}
	return nil
}

// SaveYAML 保存YAML数据到文件
func SaveYAML(filePath string, v interface{}) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("serialization of YAML data failed: %v", err)
	}

	if err := SaveToFile(filePath, data); err != nil {
		return fmt.Errorf("failed to save the YAML file: %v", err)
	}

	return nil
}
