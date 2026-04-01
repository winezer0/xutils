package cmdutils

import "testing"

// TestParseKeyValuesRules_Standard 纯标准库测试版本
func TestParseKeyValuesRules_Standard(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "正常单key解析",
			input: []string{
				"name=list:aaa,bbb",
			},
			want: map[string][]string{
				"name": {"aaa", "bbb"},
			},
		},
		{
			name: "相同key合并",
			input: []string{
				"fruit=list:apple,banana",
				"fruit=list:orange,apple",
			},
			want: map[string][]string{
				"fruit": {"apple", "banana", "orange"},
			},
		},
		{
			name: "格式错误",
			input: []string{
				"no-equal-sign",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules, err := ParseKeyValuesRules(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("err:%v wantErr:%v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			got := make(map[string]map[string]bool)
			for _, r := range rules {
				got[r.Key] = r.Values
			}

			// 对比预期
			for k, wantVals := range tt.want {
				valMap, ok := got[k]
				if !ok {
					t.Fatalf("key %s not found", k)
				}
				for _, v := range wantVals {
					if !valMap[v] {
						t.Errorf("key %s missing value %s", k, v)
					}
				}
			}
		})
	}
}
