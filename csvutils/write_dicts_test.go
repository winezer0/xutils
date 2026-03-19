package csvutils

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"
)

// 测试文件路径常量
const (
	testCSVFile  = "test_write_dicts.csv"
	testReadFile = "test_read_header.csv"
)

// 测试前清理文件
func cleanupTestFile(t *testing.T, filePath string) {
	t.Helper()
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		t.Fatalf("cleanup test file failed: %v", err)
	}
}

// 读取CSV文件所有行，用于校验测试结果（原生断言版）
func readCSVAllRows(t *testing.T, filePath string, delimiter rune) [][]string {
	t.Helper()
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("open test file failed: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter

	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("read csv rows failed: %v", err)
	}
	return rows
}

// ------------------- WriteDictsToCSV 测试用例（原生断言） -------------------

// TestWriteDictsToCSV_Normal 测试正常写入（覆盖模式+指定表头）
func TestWriteDictsToCSV_Normal(t *testing.T) {
	// 前置清理
	cleanupTestFile(t, testCSVFile)
	defer cleanupTestFile(t, testCSVFile)

	// 测试数据
	dicts := []map[string]interface{}{
		{"name": "张三", "age": 25, "city": "北京"},
		{"name": "李四", "age": 30, "city": "上海"},
	}
	header := []string{"name", "age", "city"}
	delimiter := ','
	overwrite := true

	// 执行写入
	err := WriteDictsToCSV(testCSVFile, dicts, header, delimiter, overwrite)
	if err != nil {
		t.Fatalf("WriteDictsToCSV normal write failed: %v", err)
	}

	// 校验结果
	rows := readCSVAllRows(t, testCSVFile, delimiter)
	// 校验行数：表头1行 + 数据2行 = 3行
	if len(rows) != 3 {
		t.Errorf("csv row count mismatch: expected 3, got %d", len(rows))
	}
	// 校验表头
	if !reflect.DeepEqual(rows[0], header) {
		t.Errorf("csv header mismatch: expected %v, got %v", header, rows[0])
	}
	// 校验数据行
	if !reflect.DeepEqual(rows[1], []string{"张三", "25", "北京"}) {
		t.Errorf("first data row mismatch: expected %v, got %v", []string{"张三", "25", "北京"}, rows[1])
	}
	if !reflect.DeepEqual(rows[2], []string{"李四", "30", "上海"}) {
		t.Errorf("second data row mismatch: expected %v, got %v", []string{"李四", "30", "上海"}, rows[2])
	}
}

// TestWriteDictsToCSV_AutoHeader 测试自动生成表头（未指定header）
func TestWriteDictsToCSV_AutoHeader(t *testing.T) {
	cleanupTestFile(t, testCSVFile)
	defer cleanupTestFile(t, testCSVFile)

	// 测试数据（字典键顺序随机，验证自动排序）
	dicts := []map[string]interface{}{
		{"city": "广州", "age": 28, "name": "王五"}, // 键顺序打乱
	}
	delimiter := ','
	overwrite := true

	// 执行写入
	err := WriteDictsToCSV(testCSVFile, dicts, nil, delimiter, overwrite)
	if err != nil {
		t.Fatalf("WriteDictsToCSV auto header failed: %v", err)
	}

	// 校验结果（自动生成的表头应排序为 age, city, name）
	rows := readCSVAllRows(t, testCSVFile, delimiter)
	if len(rows) != 2 {
		t.Errorf("auto header row count mismatch: expected 2, got %d", len(rows))
	}
	// 校验自动生成的表头（排序后）
	expectedHeader := []string{"age", "city", "name"}
	if !reflect.DeepEqual(rows[0], expectedHeader) {
		t.Errorf("auto header mismatch: expected %v, got %v", expectedHeader, rows[0])
	}
	// 校验数据行
	expectedRow := []string{"28", "广州", "王五"}
	if !reflect.DeepEqual(rows[1], expectedRow) {
		t.Errorf("auto header data row mismatch: expected %v, got %v", expectedRow, rows[1])
	}
}

// TestWriteDictsToCSV_Append 测试追加模式（表头匹配）
func TestWriteDictsToCSV_Append(t *testing.T) {
	cleanupTestFile(t, testCSVFile)
	defer cleanupTestFile(t, testCSVFile)

	// 第一步：覆盖模式写入初始数据
	initDicts := []map[string]interface{}{
		{"name": "张三", "age": 25},
	}
	header := []string{"name", "age"}
	err := WriteDictsToCSV(testCSVFile, initDicts, header, ',', true)
	if err != nil {
		t.Fatalf("init write failed: %v", err)
	}

	// 第二步：追加模式写入新数据
	appendDicts := []map[string]interface{}{
		{"name": "李四", "age": 30},
	}
	err = WriteDictsToCSV(testCSVFile, appendDicts, header, ',', false)
	if err != nil {
		t.Fatalf("append write failed: %v", err)
	}

	// 校验结果：表头1行 + 初始1行 + 追加1行 = 3行
	rows := readCSVAllRows(t, testCSVFile, ',')
	if len(rows) != 3 {
		t.Errorf("append row count mismatch: expected 3, got %d", len(rows))
	}
	// 校验追加数据行
	expectedRow := []string{"李四", "30"}
	if !reflect.DeepEqual(rows[2], expectedRow) {
		t.Errorf("append data row mismatch: expected %v, got %v", expectedRow, rows[2])
	}
}

// TestWriteDictsToCSV_AppendHeaderMismatch 测试追加模式（表头不匹配，写入新表头+警告）
func TestWriteDictsToCSV_AppendHeaderMismatch(t *testing.T) {
	cleanupTestFile(t, testCSVFile)
	defer cleanupTestFile(t, testCSVFile)

	// 第一步：写入初始表头 [name, age]
	initDicts := []map[string]interface{}{{"name": "张三", "age": 25}}
	err := WriteDictsToCSV(testCSVFile, initDicts, []string{"name", "age"}, ',', true)
	if err != nil {
		t.Fatalf("init write failed: %v", err)
	}

	// 第二步：追加时使用新表头 [name, age, city]（不匹配）
	newHeader := []string{"name", "age", "city"}
	appendDicts := []map[string]interface{}{{"name": "李四", "age": 30, "city": "上海"}}
	err = WriteDictsToCSV(testCSVFile, appendDicts, newHeader, ',', false)
	if err != nil {
		t.Fatalf("append header mismatch write failed: %v", err)
	}

	// 校验结果：会写入新表头，最终行数=初始2行 + 新表头1行 + 新数据1行 = 4行
	rows := readCSVAllRows(t, testCSVFile, ',')
	if len(rows) != 4 {
		t.Errorf("header mismatch append row count mismatch: expected 4, got %d", len(rows))
	}
	// 校验新表头和数据
	if !reflect.DeepEqual(rows[2], newHeader) {
		t.Errorf("new header mismatch: expected %v, got %v", newHeader, rows[2])
	}
	expectedRow := []string{"李四", "30", "上海"}
	if !reflect.DeepEqual(rows[3], expectedRow) {
		t.Errorf("header mismatch data row mismatch: expected %v, got %v", expectedRow, rows[3])
	}
}

// TestWriteDictsToCSV_EmptyDicts 测试空字典列表（不创建文件）
func TestWriteDictsToCSV_EmptyDicts(t *testing.T) {
	cleanupTestFile(t, testCSVFile)
	defer cleanupTestFile(t, testCSVFile)

	// 空字典列表
	dicts := []map[string]interface{}{}
	header := []string{"name", "age"}

	// 执行写入
	err := WriteDictsToCSV(testCSVFile, dicts, header, ',', true)
	if err != nil {
		t.Fatalf("empty dicts write failed: %v", err)
	}

	// 校验：文件不应存在
	_, err = os.Stat(testCSVFile)
	if !os.IsNotExist(err) {
		t.Errorf("empty dicts should not create file, but file exists")
	}
}

// TestWriteDictsToCSV_CustomDelimiter 测试自定义分隔符（分号）
func TestWriteDictsToCSV_CustomDelimiter(t *testing.T) {
	cleanupTestFile(t, testCSVFile)
	defer cleanupTestFile(t, testCSVFile)

	// 测试数据
	dicts := []map[string]interface{}{
		{"name": "赵六", "age": 35},
	}
	header := []string{"name", "age"}
	delimiter := ';'

	// 执行写入
	err := WriteDictsToCSV(testCSVFile, dicts, header, delimiter, true)
	if err != nil {
		t.Fatalf("custom delimiter write failed: %v", err)
	}

	// 校验结果
	rows := readCSVAllRows(t, testCSVFile, delimiter)
	if len(rows) != 2 {
		t.Errorf("custom delimiter row count mismatch: expected 2, got %d", len(rows))
	}
	expectedRow := []string{"赵六", "35"}
	if !reflect.DeepEqual(rows[1], expectedRow) {
		t.Errorf("custom delimiter data row mismatch: expected %v, got %v", expectedRow, rows[1])
	}
}

// TestWriteDictsToCSV_EmptyDelimiter 测试空分隔符（默认使用逗号）
func TestWriteDictsToCSV_EmptyDelimiter(t *testing.T) {
	cleanupTestFile(t, testCSVFile)
	defer cleanupTestFile(t, testCSVFile)

	// 空分隔符（rune 0）
	dicts := []map[string]interface{}{{"name": "钱七", "age": 40}}
	header := []string{"name", "age"}
	delimiter := rune(0) // 空分隔符

	// 执行写入
	err := WriteDictsToCSV(testCSVFile, dicts, header, delimiter, true)
	if err != nil {
		t.Fatalf("empty delimiter write failed: %v", err)
	}

	// 校验：默认使用逗号分隔
	rows := readCSVAllRows(t, testCSVFile, ',')
	if len(rows) != 2 {
		t.Errorf("empty delimiter row count mismatch: expected 2, got %d", len(rows))
	}
	expectedRow := []string{"钱七", "40"}
	if !reflect.DeepEqual(rows[1], expectedRow) {
		t.Errorf("empty delimiter data row mismatch: expected %v, got %v", expectedRow, rows[1])
	}
}
