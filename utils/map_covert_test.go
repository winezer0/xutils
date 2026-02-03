package utils

import (
	"testing"
)

func TestGetMapBool(t *testing.T) {
	tests := []struct {
		name     string
		inputMap map[string]interface{}
		keys     []string
		want     bool
	}{
		{
			name: "Field exists with true value",
			inputMap: map[string]interface{}{
				"tracerCode": true,
			},
			keys: []string{"tracerCode"},
			want: true,
		},
		{
			name: "Field exists with false value",
			inputMap: map[string]interface{}{
				"tracerCode": false,
			},
			keys: []string{"tracerCode"},
			want: false,
		},
		{
			name: "Field exists with alias true value",
			inputMap: map[string]interface{}{
				"tracer_code": true,
			},
			keys: []string{"tracerCode", "tracer_code"},
			want: true,
		},
		{
			name: "Field not exists",
			inputMap: map[string]interface{}{
				"otherField": true,
			},
			keys: []string{"tracerCode"},
			want: false,
		},
		{
			name: "Field value is not bool",
			inputMap: map[string]interface{}{
				"tracerCode": "true",
			},
			keys: []string{"tracerCode"},
			want: false,
		},
		{
			name: "Field value is nil",
			inputMap: map[string]interface{}{
				"tracerCode": nil,
			},
			keys: []string{"tracerCode"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMapBool(tt.inputMap, tt.keys...)
			if got != tt.want {
				t.Errorf("GetMapBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMapStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		inputMap map[string]interface{}
		keys     []string
		want     []string
	}{
		{
			name: "Field exists with []string type",
			inputMap: map[string]interface{}{
				"dependFuncs": []string{"func1", "func2", "func3"},
			},
			keys: []string{"dependFuncs"},
			want: []string{"func1", "func2", "func3"},
		},
		{
			name: "Field exists with []interface{} type containing strings",
			inputMap: map[string]interface{}{
				"dependFuncs": []interface{}{"func1", "func2", "func3"},
			},
			keys: []string{"dependFuncs"},
			want: []string{"func1", "func2", "func3"},
		},
		{
			name: "Field exists with []interface{} type containing mixed types",
			inputMap: map[string]interface{}{
				"dependFuncs": []interface{}{"func1", 123, true, "func2"},
			},
			keys: []string{"dependFuncs"},
			want: []string{"func1", "func2"},
		},
		{
			name: "Field exists with alias containing []string",
			inputMap: map[string]interface{}{
				"functionList": []string{"func1", "func2"},
			},
			keys: []string{"dependFuncs", "functionList"},
			want: []string{"func1", "func2"},
		},
		{
			name: "Field exists with alias containing []interface{}",
			inputMap: map[string]interface{}{
				"function_list": []interface{}{"func1", "func2"},
			},
			keys: []string{"dependFuncs", "function_list"},
			want: []string{"func1", "func2"},
		},
		{
			name: "Field not exists",
			inputMap: map[string]interface{}{
				"otherField": []string{"func1", "func2"},
			},
			keys: []string{"dependFuncs"},
			want: nil,
		},
		{
			name: "Field value is not a slice",
			inputMap: map[string]interface{}{
				"dependFuncs": "func1,func2,func3",
			},
			keys: []string{"dependFuncs"},
			want: nil,
		},
		{
			name: "Field value is nil",
			inputMap: map[string]interface{}{
				"dependFuncs": nil,
			},
			keys: []string{"dependFuncs"},
			want: nil,
		},
		{
			name: "Empty slice",
			inputMap: map[string]interface{}{
				"dependFuncs": []string{},
			},
			keys: []string{"dependFuncs"},
			want: []string{},
		},
		{
			name: "Empty interface slice",
			inputMap: map[string]interface{}{
				"dependFuncs": []interface{}{},
			},
			keys: []string{"dependFuncs"},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMapStringSlice(tt.inputMap, tt.keys...)
			
			// 检查切片长度
			if len(got) != len(tt.want) {
				t.Errorf("GetMapStringSlice() length = %d, want %d", len(got), len(tt.want))
				return
			}
			
			// 检查切片内容
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetMapStringSlice() index %d = %s, want %s", i, v, tt.want[i])
				}
			}
		})
	}
}
