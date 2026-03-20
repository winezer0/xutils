package csvutils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ReadCSVToStructs 将 CSV 文件内容映射到结构体切片
// 参数:
//
//	csvFile: CSV 文件路径
//	out: 必须是指向结构体切片的指针（例如 &[]User{}）
//	tagName: 结构体字段标签名（例如 "csv"），用于匹配 CSV 表头
func ReadCSVToStructs(csvFile string, out interface{}, tagName string) error {
	// 1. 先调用原有函数读取 CSV 为字典格式
	_, dicts, err := ReadCSVToDicts(csvFile)
	if err != nil {
		return fmt.Errorf("failed to read CSV to dicts: %w", err)
	}

	// 2. 验证 out 参数的类型（必须是指向结构体切片的指针）
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr || outVal.IsNil() {
		return fmt.Errorf("out must be a non-nil pointer to a slice of structs")
	}

	sliceVal := outVal.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("out must point to a slice (got %s)", sliceVal.Kind())
	}

	// 获取结构体的类型（切片的元素类型）
	elemType := sliceVal.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("slice element must be a struct (got %s)", elemType.Kind())
	}

	// 3. 构建表头到结构体字段的映射（支持 tag 匹配）
	headerToField := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		// 获取字段的 tag 值（例如 `csv:"name"` 中的 "name"）
		tagValue := field.Tag.Get(tagName)
		if tagValue == "" {
			// 如果没有 tag，使用字段名（不区分大小写）
			tagValue = field.Name
		}
		headerToField[strings.ToLower(tagValue)] = i
	}

	// 4. 将每个字典转换为结构体实例
	for dictIdx, dict := range dicts {
		// 创建新的结构体实例
		structVal := reflect.New(elemType).Elem()

		// 遍历表头，给结构体字段赋值
		for header, value := range dict {
			headerLower := strings.ToLower(header)
			fieldIdx, ok := headerToField[headerLower]
			if !ok {
				// 忽略不存在的字段（也可以选择返回错误）
				continue
			}

			field := structVal.Field(fieldIdx)
			if !field.CanSet() {
				continue // 跳过不可设置的字段（例如私有字段）
			}

			// 根据字段类型转换值
			err := setFieldValue(field, value)
			if err != nil {
				return fmt.Errorf("row %d, field %s: %w", dictIdx+2, header, err)
			}
		}

		// 将结构体添加到切片
		sliceVal.Set(reflect.Append(sliceVal, structVal))
	}

	return nil
}

// setFieldValue 根据字段类型设置值（支持常见类型）
func setFieldValue(field reflect.Value, value string) error {
	if value == "" {
		return nil // 空值不设置
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int value: %w", err)
		}
		field.SetInt(val)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid uint value: %w", err)
		}
		field.SetUint(val)

	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %w", err)
		}
		field.SetFloat(val)

	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid bool value: %w", err)
		}
		field.SetBool(val)

	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
