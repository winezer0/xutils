package csvutils

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type User struct {
	Name  string `csv:"name" json:"name" yaml:"user_name"`
	Age   int    `csv:"age" json:"age" yaml:"user_age"`
	Email string `csv:"email" json:"email" yaml:"user_email"`
}

type Product struct {
	ID       int     `csv:"id" json:"id" yaml:"product_id"`
	Name     string  `csv:"name" json:"name" yaml:"product_name"`
	Price    float64 `csv:"price" json:"price" yaml:"product_price"`
	Quantity int     `csv:"quantity" json:"quantity" yaml:"product_quantity"`
}

func TestReadCSVToStructs(t *testing.T) {
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "test.csv")

	csvContent := `name,age,email
张三,25,zhangsan@example.com
李四,30,lisi@example.com
王五,28,wangwu@example.com`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var users []User
	err = ReadCSVToStructs(csvFile, &users, "csv")
	if err != nil {
		t.Fatalf("读取CSV到结构体失败: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("期望3个用户，实际%d个", len(users))
	}

	expected := User{Name: "张三", Age: 25, Email: "zhangsan@example.com"}
	if users[0].Name != expected.Name || users[0].Age != expected.Age || users[0].Email != expected.Email {
		t.Errorf("第一个用户数据不匹配: 期望 %+v，实际 %+v", expected, users[0])
	}

	expectedAge, _ := strconv.ParseInt("30", 10, 64)
	if users[1].Age != int(expectedAge) {
		t.Errorf("第二个用户年龄不匹配: 期望 30，实际 %d", users[1].Age)
	}
}

func TestReadCSVToStructsWithDifferentTags(t *testing.T) {
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "test2.csv")

	csvContent := `id,name,price,quantity
1,苹果,3.5,100
2,香蕉,2.0,150
3,橙子,4.0,80`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var products []Product
	err = ReadCSVToStructs(csvFile, &products, "csv")
	if err != nil {
		t.Fatalf("读取CSV到结构体失败: %v", err)
	}

	if len(products) != 3 {
		t.Errorf("期望3个产品，实际%d个", len(products))
	}

	expected := Product{ID: 1, Name: "苹果", Price: 3.5, Quantity: 100}
	if products[0].ID != expected.ID || products[0].Name != expected.Name || products[0].Price != expected.Price || products[0].Quantity != expected.Quantity {
		t.Errorf("第一个产品数据不匹配: 期望 %+v，实际 %+v", expected, products[0])
	}
}

func TestReadCSVToStructsEmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "empty.csv")

	csvContent := `name,age,email`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var users []User
	err = ReadCSVToStructs(csvFile, &users, "csv")
	if err != nil {
		t.Fatalf("读取空CSV文件失败: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("期望0个用户，实际%d个", len(users))
	}
}

func TestReadCSVToStructsMissingField(t *testing.T) {
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "missing.csv")

	csvContent := `name,email
张三,zhangsan@example.com
李四,lisi@example.com`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var users []User
	err = ReadCSVToStructs(csvFile, &users, "csv")
	if err != nil {
		t.Fatalf("读取CSV到结构体失败: %v", err)
	}

	if users[0].Age != 0 {
		t.Errorf("缺失字段应该为零值，实际为 %d", users[0].Age)
	}
}

func TestReadCSVToStructsInvalidType(t *testing.T) {
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "invalid.csv")

	csvContent := `name,age,email
张三,abc,zhangsan@example.com`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var users []User
	err = ReadCSVToStructs(csvFile, &users, "csv")
	if err == nil {
		t.Error("期望错误，但没有返回错误")
	}
}

func TestReadCSVToStructsSimple(t *testing.T) {
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "simple.csv")

	csvContent := `name,age,email
张三,25,zhangsan@example.com
李四,30,lisi@example.com`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var users []User
	err = ReadCSVToStructsSimple(csvFile, &users)
	if err != nil {
		t.Fatalf("使用简单方法读取CSV到结构体失败: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("期望2个用户，实际%d个", len(users))
	}

	expected := User{Name: "张三", Age: 25, Email: "zhangsan@example.com"}
	if users[0].Name != expected.Name || users[0].Age != expected.Age || users[0].Email != expected.Email {
		t.Errorf("第一个用户数据不匹配: 期望 %+v，实际 %+v", expected, users[0])
	}
}

func TestReadCSVToStructsSimpleWithProducts(t *testing.T) {
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "products.csv")

	csvContent := `id,name,price,quantity
1,苹果,3.5,100
2,香蕉,2.0,150
3,橙子,4.0,80`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var products []Product
	err = ReadCSVToStructsSimple(csvFile, &products)
	if err != nil {
		t.Fatalf("使用简单方法读取CSV到结构体失败: %v", err)
	}

	if len(products) != 3 {
		t.Errorf("期望3个产品，实际%d个", len(products))
	}

	expected := Product{ID: 1, Name: "苹果", Price: 3.5, Quantity: 100}
	if products[0].ID != expected.ID || products[0].Name != expected.Name || products[0].Price != expected.Price || products[0].Quantity != expected.Quantity {
		t.Errorf("第一个产品数据不匹配: 期望 %+v，实际 %+v", expected, products[0])
	}
}

func TestReadCSVToStructsSimpleEmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "empty2.csv")

	csvContent := `name,age,email`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var users []User
	err = ReadCSVToStructsSimple(csvFile, &users)
	if err != nil {
		t.Fatalf("使用简单方法读取空CSV文件失败: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("期望0个用户，实际%d个", len(users))
	}
}

func TestReadCSVToStructsSimpleInvalidType(t *testing.T) {
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "invalid2.csv")

	csvContent := `name,age,email
张三,abc,zhangsan@example.com`

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var users []User
	err = ReadCSVToStructsSimple(csvFile, &users)
	if err == nil {
		t.Error("期望错误，但没有返回错误")
	}
}
