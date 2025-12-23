package kenv

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const envFileName = ".env"

// LoadDotEnv function    加载 .env 文件到环境变量.
func LoadDotEnv() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("get exec dir failed: %w", err)
	}

	file := filepath.Join(dir, envFileName)
	if err = loadEnvFile(file); err != nil {
		// 失败后尝试当前目录
		if err = loadEnvFile(envFileName); err != nil {
			if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
				fmt.Printf("[TWIG] load [ %s ] failed: %v\n", envFileName, err)
				return nil
			}
			return fmt.Errorf("load [ %s ] failed: %w", envFileName, err)
		}
	}

	fmt.Printf("[TWIG] load [ %s ] success\n", envFileName)
	return nil
}

// LoadEnvFile function    加载指定的环境变量文件.
func LoadEnvFile(filename string) error {
	return loadEnvFile(filename)
}

func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 移除引号
		value = strings.Trim(value, `"'`)

		// 设置环境变量(不覆盖已存在的)
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
