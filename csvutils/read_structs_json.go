package csvutils

import (
	"fmt"
	"reflect"
)

// ReadCSVToStructsSimple 基于CSV转JSON的方式，将CSV读取为结构体切片 CSV → JSON → Struct
// target 需要传入结构体切片的指针（如 &[]User{}）
func ReadCSVToStructsSimple(csvFile string, target interface{}) error {
	// 1. 先用原有函数读取为map数组
	_, dicts, err := ReadCSVToDicts(csvFile)
	if err != nil {
		return fmt.Errorf("read CSV to dicts failed: %w", err)
	}

	// 2. 验证 target 参数的类型
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer to a slice of structs")
	}

	sliceVal := targetVal.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("target must point to a slice (got %s)", sliceVal.Kind())
	}

	elemType := sliceVal.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("slice element must be a struct (got %s)", elemType.Kind())
	}

	// 3. 构建字段映射
	fieldMap := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tagValue := field.Tag.Get("json")
		if tagValue == "" {
			tagValue = field.Name
		}
		fieldMap[tagValue] = i
	}

	// 4. 创建目标切片
	resultSlice := reflect.MakeSlice(sliceVal.Type(), 0, len(dicts))

	for _, dict := range dicts {
		structVal := reflect.New(elemType).Elem()

		for header, value := range dict {
			fieldIdx, ok := fieldMap[header]
			if !ok {
				continue
			}

			field := structVal.Field(fieldIdx)
			if !field.CanSet() {
				continue
			}

			err := setFieldValue(field, value)
			if err != nil {
				return fmt.Errorf("field %s: %w", header, err)
			}
		}

		resultSlice = reflect.Append(resultSlice, structVal)
	}

	sliceVal.Set(resultSlice)
	return nil
}
