package loader

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/spelens-gud/twig/kits/kgo"
	"github.com/spelens-gud/twig/version"
)

var (
	// execPath 执行文件路径
	execPath = os.Args[0]
	// execAbsoluteDir 执行文件相对路径
	execAbsoluteDir = func() string {
		ret, err := filepath.Abs(filepath.Dir(execPath))
		if err != nil {
			panic(err)
		}
		return ret
	}()
	// appname 应用名称
	appName = NewConfigSource(getApplicationAppName)
	// namepace 命名空间
	namespace = NewConfigSource(getApplicationNameSpace)
	// hostname // 主机名
	hostname = NewConfigSource(getHostname)
	// execMd5
	execMd5 = NewConfigSource(getExecMd5)
	// 获取环境服务版本
	env = NewConfigSource(func() string {
		ret := kgo.GetFirstValidString(os.Getenv(EnvKey), Development)
		fmt.Printf("[SY-TWIG] env loaded [ %s ] version [ %s ]\n", ret, version.Get().Short())
		return ret
	})
)

// GetEnv 获取环境环境版本.
func GetEnv() Env {
	return Env(env.Get())
}

// SetEnv 设置环境版本.
func SetEnv(e string) {
	fmt.Printf("[SY-TWIG] env updated [ %s ]\n", e)
	env.Set(e)
}

// ExecDir 获取执行文件目录.
func ExecDir() string { return execAbsoluteDir }

// ExecName 获取执行文件名.
func ExecName() string { return filepath.Base(execPath) }

// HostName 获取Hostname.
func HostName() string { return hostname.Get() }

// GetLogFileDir 获取日志输出目录.
func GetLogFileDir() string {
	if dir := os.Getenv(EnvKeyLogFileDir); len(dir) > 0 {
		return dir
	}
	return ""
}

// SetLogFileDir 设置日志输出目录.
func SetLogFileDir(dir string) {
	_ = os.Setenv(EnvKeyLogFileDir, dir)
}

// GetApplicationNameSpace 获取应用命名空间.
func GetApplicationNameSpace() string { return namespace.Get() }

// SetApplicationNameSpace 设置应用命名空间.
func SetApplicationNameSpace(str string) { namespace.Set(str) }

// GetApplicationAppName 获取应用名.
func GetApplicationAppName() string { return appName.Get() }

// SetApplicationAppName 设置应用名.
func SetApplicationAppName(str string) { appName.Set(str) }

// GetApplicationName 获取应用完整名 格式: $namespace:$app.
func GetApplicationName() string { return GetApplicationNameSpace() + ":" + GetApplicationAppName() }

// GetApplicationServiceName 获取应用服务名 格式: $app.$namespace.
func GetApplicationServiceName() string {
	return GetApplicationAppName() + "." + GetApplicationNameSpace()
}

// initFormat 格式化连接字符
func initFormat(in string) string {
	return strings.ReplaceAll(strings.ToLower(in), "_", "-")
}

// getHostname 获取主机名称
func getHostname() (n string) { n, _ = os.Hostname(); return n }

// getApplicationNameSpace 获取应用命名空间
func getApplicationNameSpace() (n string) {
	if n = os.Getenv(EnvKeyContainerNamespace); len(n) == 0 {
		n = "unknown"
	}
	return initFormat(n)
}

// getApplicationAppName // 获取应用名称
func getApplicationAppName() (n string) {
	defaultFileName := "unknown-app-"
	if hash := execMd5.Get(); len(hash) >= 8 {
		defaultFileName += hash[:8]
	} else {
		defaultFileName += uuid.New().String()[:8]
	}
	return initFormat(kgo.GetFirstValidString(os.Getenv(EnvKeyContainerAppName), defaultFileName))
}

// getExecMd5 获取执行文件的md5
func getExecMd5() (n string) {
	md5Hash, _ := getFileMd5(execPath)
	return md5Hash
}

// getFileMd5 获取文件的md5
func getFileMd5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	//nolint:errcheck
	defer file.Close()
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
