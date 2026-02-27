package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCheckFilesExist 核心测试函数：覆盖全场景
func TestCheckFilesExist(t *testing.T) {
	// ========== 步骤1：创建测试用临时文件/目录（自动清理） ==========
	testDir := filepath.Join(os.TempDir(), "test_check_files")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Logf("清理测试目录失败: %v", err)
		}
	}()

	// 1. 创建存在的文件
	existFile := filepath.Join(testDir, "exist.txt")
	if err := os.WriteFile(existFile, []byte("test"), 0644); err != nil {
		t.Fatalf("创建存在的文件失败: %v", err)
	}

	// 2. 创建存在的目录
	existDir := filepath.Join(testDir, "exist_dir")
	if err := os.Mkdir(existDir, 0755); err != nil {
		t.Fatalf("创建存在的目录失败: %v", err)
	}

	// 3. 不存在的文件路径
	nonExistFile := filepath.Join(testDir, "non_exist.txt")

	// 4. 本地存在的文件 (当前工作目录)
	localExistFile := "local_exist_test.txt"
	if err := os.WriteFile(localExistFile, []byte("test"), 0644); err != nil {
		t.Fatalf("创建本地存在文件失败: %v", err)
	}
	defer os.Remove(localExistFile)

	// 5. Windows 混合分隔符路径 (模拟真实场景，Go 通常能自动处理，但测试一下无妨)
	mixedSepFile := strings.ReplaceAll(existFile, "\\", "/")

	// ========== 步骤2：定义测试用例 ==========
	type testCase struct {
		name       string
		inputPaths []string
		wantValid  []string // 预期存在的路径
		wantPlain  []string // 预期不存在/纯字符串
	}

	testCases := []testCase{
		{
			name: "normal_case",
			inputPaths: []string{
				existFile,      // 存在的文件
				existDir,       // 存在的目录 (os.Stat 对目录也返回 nil)
				nonExistFile,   // 不存在的文件
				localExistFile, // 本地存在的文件
				"hello world",  // 纯字符串
				"",             // 空字符串
				mixedSepFile,   // 混合分隔符存在文件
			},
			wantValid: []string{
				existFile,
				existDir,
				localExistFile,
				mixedSepFile,
			},
			wantPlain: []string{
				nonExistFile,
				"hello world",
				"",
			},
		},
		{
			name:       "empty_input",
			inputPaths: []string{},
			wantValid:  []string{},
			wantPlain:  []string{},
		},
		{
			name: "all_non_exist",
			inputPaths: []string{
				"non_exist_1.txt",
				"test@file.txt",
				"a|b",
				"   ", // 全空格字符串
			},
			wantValid: []string{},
			wantPlain: []string{
				"non_exist_1.txt",
				"test@file.txt",
				"a|b",
				"   ",
			},
		},
		{
			name: "error_case", // 异常路径（如超长路径）
			inputPaths: []string{
				filepath.Join(testDir, strings.Repeat("a", 300)+".txt"),
			},
			wantValid: []string{},
			wantPlain: []string{
				filepath.Join(testDir, strings.Repeat("a", 300)+".txt"),
			},
		},
	}

	// ========== 步骤3：执行测试 ==========
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 【关键修改】接收两个返回值，而不是结构体
			validFilePaths, plainStrings := CheckPathsExist(tc.inputPaths)

			// 校验存在的文件路径
			if !sliceEqual(validFilePaths, tc.wantValid) {
				t.Errorf("ValidFilePaths 不匹配:\ngot  %v\nwant %v", validFilePaths, tc.wantValid)
			}

			// 校验不存在/纯字符串
			if !sliceEqual(plainStrings, tc.wantPlain) {
				t.Errorf("PlainStrings 不匹配:\ngot  %v\nwant %v", plainStrings, tc.wantPlain)
			}
		})
	}
}

// sliceEqual 辅助函数：严格比较字符串切片（顺序 + 内容完全一致）
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
