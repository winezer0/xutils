package cmdutils

import (
	"os"
	"testing"
)

// TestParseRule 单元测试：测试单条规则解析函数 parseRule
// 覆盖：list模式、默认模式、去重去空、异常格式、文件读取
func TestParseRule(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name    string   // 用例名称
		input   string   // 输入
		want    []string // 期望输出
		wantErr bool     // 是否期望错误
	}{
		{
			name:  "list前缀正常解析",
			input: "list:aaa,bbb,ccc",
			want:  []string{"aaa", "bbb", "ccc"},
		},
		{
			name:  "无前缀默认list解析",
			input: "xxx,yyy,zzz",
			want:  []string{"xxx", "yyy", "zzz"},
		},
		{
			name:  "带空格自动清理",
			input: "list:  aaa ,  bbb  ,,,ccc",
			want:  []string{"aaa", "bbb", "ccc"},
		},
		{
			name:  "内部重复值自动去重",
			input: "list:aaa,aaa,bbb,bbb",
			want:  []string{"aaa", "bbb"},
		},
		{
			name:  "空内容返回空列表",
			input: "list:",
			want:  []string{},
		},
		{
			name:    "文件不存在返回错误",
			input:   "file:not_exist_file_1234.txt",
			wantErr: true,
		},
	}

	// 执行所有测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseValueRule(tt.input)

			// 判断错误是否符合预期
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseRule() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			// 比较结果是否一致
			if !sliceEqual(got, tt.want) {
				t.Errorf("parseRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseRule_File 专门测试 file: 前缀读取真实文件
func TestParseRule_File(t *testing.T) {
	// 1. 创建临时测试文件
	content := "line1\n\nline2\n  line3  \nline1\n"
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // 测试完删除文件

	// 2. 写入测试内容
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// 3. 调用 parseRule 解析 file: 规则
	got, err := ParseValueRule("file:" + tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 4. 期望结果：去空、去重、去空格
	want := []string{"line1", "line2", "line3"}
	if !sliceEqual(got, want) {
		t.Errorf("file parse got = %v, want %v", got, want)
	}
}

// sliceEqual 辅助函数：判断两个字符串切片是否完全相等
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
