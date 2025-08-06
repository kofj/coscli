package util

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetContent 获取配置信息（若是本地文件目录则读取文件内容）
func GetContent(input string) ([]byte, error) {
	// 处理文件路径
	if strings.HasPrefix(input, "file://") {
		filePath := strings.TrimPrefix(input, "file://")
		return readFile(filePath)
	}

	// 不是文件路径则直接返回输入内容
	return []byte(input), nil
}

// readFile 读取本地文件目录文件内容
func readFile(filePath string) ([]byte, error) {
	// 解析绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// 检查路径是否存在
	fileInfo, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", absPath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}

	// 检查是否是目录
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file: %s", absPath)
	}

	// 检查文件大小
	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("file is empty: %s", absPath)
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 检查内容是否为空
	if len(content) == 0 {
		return nil, fmt.Errorf("file contains no data: %s", absPath)
	}

	logger.Info("Loaded configuration from file: %s (%d bytes)", absPath, len(content))
	return content, nil
}

func ParseContent[T any](content []byte, target *T) error {
	// 检查内容是否为空
	if len(content) == 0 {
		return fmt.Errorf("content is empty")
	}

	// 尝试解析为JSON
	if err := json.Unmarshal(content, &target); err == nil {
		logger.Info("Detected JSON format")
		return nil
	}

	// 尝试解析为XML
	if err := xml.Unmarshal(content, &target); err == nil {
		logger.Info("Detected XML format")
		return nil
	}

	// 无法识别格式
	return fmt.Errorf("unrecognized configuration format, must be JSON or XML")
}
