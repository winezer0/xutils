package utils

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"os"
)

// ReadFileWithEncoding 读取指定编码的文件内容并转换为UTF-8字符串
// 如果encoding为空或不匹配任何已知编码，则默认以UTF-8读取
// 读取时会忽略解码错误，类似Python的errors=ignore
type ignoreErrorTransformer struct {
	transform.Transformer
}

// Transform 实现Transform方法，忽略解码错误
func (t ignoreErrorTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	nDst, nSrc, err = t.Transformer.Transform(dst, src, atEOF)
	if err != nil {
		// 忽略错误，继续处理
		return nDst, nSrc, nil
	}
	return nDst, nSrc, err
}

// 辅助函数：创建忽略错误的转换器
func ignoreErrors(t transform.Transformer) transform.Transformer {
	return ignoreErrorTransformer{t}
}

// readFileWithEncoding 使用指定的编码器读取文件并解码为UTF-8字符串
// 忽略解码错误，类似Python的errors=ignore
func readFileWithEncoding(filePath string, enc encoding.Encoding) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建解码读取器，并使用自定义的错误忽略转换器
	decoder := enc.NewDecoder()
	reader := transform.NewReader(file, ignoreErrors(decoder))

	// 读取并解码内容
	content, err := io.ReadAll(reader)
	if err != nil && len(content) == 0 {
		return "", err
	}

	// 强制将字节切片转换为字符串，即使包含无效 UTF-8
	return string(content), nil
}

// ReadFileWithEncoding 读取指定编码的文件内容并转换为UTF-8字符串
// 如果encoding为空或不匹配任何已知编码，则默认以UTF-8读取
// 读取时会忽略解码错误，类似Python的errors=ignore
func ReadFileWithEncoding(filePath, encode string) (string, error) {
	encode = DetectFileEncode(filePath, encode)
	return readFileWithEncoding(filePath, NormalizedEncode(encode))
}

// ReadFileToListWithEncoding 读取文件并返回处理后的行列表
func ReadFileToListWithEncoding(filePath, encode string, deUnprint, ignoreBlanks, deWeight bool) ([]string, error) {
	if !FileExists(filePath) {
		return nil, fmt.Errorf("file is not exists")
	}
	// 分析文件编码
	encode = DetectFileEncode(filePath, encode)
	// 转换为编码类型
	enc := NormalizedEncode(encode)

	lines, strings, err2 := readFileToListWithEncoding(filePath, enc, deUnprint, ignoreBlanks)
	if err2 != nil {
		return strings, err2
	}

	// 自动去重
	if deWeight {
		lines = UniqueSlice(lines, false, ignoreBlanks)
	}

	return lines, nil
}

func readFileToListWithEncoding(filePath string, enc encoding.Encoding, deUnprint bool, ignoreBlanks bool) ([]string, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	// 关键修正：将reader声明为io.Reader接口类型(兼容所有读取器)
	var reader io.Reader = file // 默认为文件本身
	if enc != nil {
		// 应用编码转换，返回的transform.Reader实现了io.Reader接口
		reader = transform.NewReader(file, enc.NewDecoder())
	}

	var lines []string
	scanner := bufio.NewScanner(reader) // scanner接受io.Reader类型
	for scanner.Scan() {
		line := scanner.Text()
		// 清理不可打印字符
		if deUnprint {
			line = CleanupUnprintableChars(line)
		}
		// 跳过空白行
		if ignoreBlanks && line == "" {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return lines, nil, nil
}
