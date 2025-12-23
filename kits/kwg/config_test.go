package kwg

import (
	"testing"
)

func TestConfig_ParseSingleLine(t *testing.T) {
	// 单行格式配置
	content := `[Interface] PrivateKey = yPIur3SqU0Dk7gRQbbvYXxrdz8hetb/uZHsmHbj+wmU= Address = 10.0.0.12/24   [Peer] PublicKey = 3pHdc+5JxKGHo6uKZQzsRxY20zT4NH47Yy0nP9ftTzU= PresharedKey = k/gMgN/L40Q0szSXQDeNO7UzDAmjAPGeRCw9EjRHG7I= AllowedIPs = 10.0.0.0/24 PersistentKeepalive = 30 Endpoint = c2396771168ab622.natapp.cc:22331`

	cfg := &Config{}
	if err := cfg.LoadFromContent(content); err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 验证 Interface 段
	if cfg.PrivateKey != "yPIur3SqU0Dk7gRQbbvYXxrdz8hetb/uZHsmHbj+wmU=" {
		t.Errorf("PrivateKey 解析错误: %s", cfg.PrivateKey)
	}
	if len(cfg.LocalAddrs) != 1 || cfg.LocalAddrs[0] != "10.0.0.12/24" {
		t.Errorf("LocalAddrs 解析错误: %v", cfg.LocalAddrs)
	}

	// 验证 Peer 段
	if cfg.PublicKey != "3pHdc+5JxKGHo6uKZQzsRxY20zT4NH47Yy0nP9ftTzU=" {
		t.Errorf("PublicKey 解析错误: %s", cfg.PublicKey)
	}
	if cfg.PresharedKey != "k/gMgN/L40Q0szSXQDeNO7UzDAmjAPGeRCw9EjRHG7I=" {
		t.Errorf("PresharedKey 解析错误: %s", cfg.PresharedKey)
	}
	if len(cfg.AllowedIPs) != 1 || cfg.AllowedIPs[0] != "10.0.0.0/24" {
		t.Errorf("AllowedIPs 解析错误: %v", cfg.AllowedIPs)
	}
	if cfg.PersistentKeepalive != 30 {
		t.Errorf("PersistentKeepalive 解析错误: %d", cfg.PersistentKeepalive)
	}
	if cfg.Endpoint != "c2396771168ab622.natapp.cc:22331" {
		t.Errorf("Endpoint 解析错误: %s", cfg.Endpoint)
	}

	t.Logf("解析结果: %+v", cfg)
}

func TestConfig_ParseMultiLine(t *testing.T) {
	// 标准多行格式配置
	content := `[Interface]
PrivateKey = yPIur3SqU0Dk7gRQbbvYXxrdz8hetb/uZHsmHbj+wmU=
Address = 10.0.0.12/24
DNS = 8.8.8.8, 8.8.4.4

[Peer]
PublicKey = 3pHdc+5JxKGHo6uKZQzsRxY20zT4NH47Yy0nP9ftTzU=
PresharedKey = k/gMgN/L40Q0szSXQDeNO7UzDAmjAPGeRCw9EjRHG7I=
AllowedIPs = 10.0.0.0/24, 192.168.1.0/24
PersistentKeepalive = 30
Endpoint = vpn.example.com:51820`

	cfg := &Config{}
	if err := cfg.LoadFromContent(content); err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if cfg.PrivateKey != "yPIur3SqU0Dk7gRQbbvYXxrdz8hetb/uZHsmHbj+wmU=" {
		t.Errorf("PrivateKey 解析错误: %s", cfg.PrivateKey)
	}
	if len(cfg.DNSAddrs) != 2 {
		t.Errorf("DNSAddrs 解析错误: %v", cfg.DNSAddrs)
	}
	if len(cfg.AllowedIPs) != 2 {
		t.Errorf("AllowedIPs 解析错误: %v", cfg.AllowedIPs)
	}

	t.Logf("解析结果: %+v", cfg)
}

func TestConfig_ToIPC(t *testing.T) {
	cfg := &Config{
		PrivateKey:          "yPIur3SqU0Dk7gRQbbvYXxrdz8hetb/uZHsmHbj+wmU=",
		PublicKey:           "3pHdc+5JxKGHo6uKZQzsRxY20zT4NH47Yy0nP9ftTzU=",
		AllowedIPs:          []string{"10.0.0.0/24"},
		PersistentKeepalive: 30,
	}

	ipc := cfg.ToIPC()
	t.Logf("IPC 格式:\n%s", ipc)

	if ipc == "" {
		t.Error("ToIPC 返回空字符串")
	}
}
