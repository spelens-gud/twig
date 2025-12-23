package kwg

import (
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"

	"github.com/spelens-gud/twig/kits/kcrypto"
	"github.com/spelens-gud/twig/kits/kurl"
)

const (
	DefaultPersistentKeepalive = 25
)

var (
	DefaultDNSAddrs = []string{
		"223.5.5.5", "223.6.6.6",
		"119.29.29.29", "182.254.116.116",
		"180.76.76.76",
		"8.8.8.8", "8.8.4.4",
		"1.1.1.1", "1.0.0.1",
		"208.67.222.222", "208.67.220.220",
	}
)

// Config WireGuard 配置.
type Config struct {
	AllowedIPs          []string `json:"allowed_i_ps,omitempty"`         // 允许的 IP 地址
	LocalAddrs          []string `json:"local_addrs,omitempty"`          // 本地地址
	DNSAddrs            []string `json:"dns_addrs,omitempty"`            // DNS 地址
	PrivateKey          string   `json:"private_key,omitempty"`          // 私钥
	PublicKey           string   `json:"public_key,omitempty"`           // 公钥
	PresharedKey        string   `json:"preshared_key,omitempty"`        // 预共享密钥
	Endpoint            string   `json:"endpoint,omitempty"`             // 端点
	PersistentKeepalive int      `json:"persistent_keepalive,omitempty"` // 持久心跳
	LogLevel            int      `json:"log_level,omitempty"`            // 日志级别: 0=silent, 1=error, 2=verbose
}

func NewConfig() *Config {
	return &Config{}
}

// Init 初始化默认值.
func (c *Config) Init() {
	if c.PersistentKeepalive == 0 {
		c.PersistentKeepalive = DefaultPersistentKeepalive
	}
	if c.DNSAddrs == nil {
		c.DNSAddrs = DefaultDNSAddrs
	}
	if c.AllowedIPs == nil {
		c.AllowedIPs = []string{"0.0.0.0/0", "::/0"}
	}
	if c.LocalAddrs == nil {
		c.LocalAddrs = make([]string, 0)
	}
}

// LoadFromFile 从文件加载配置.
func (c *Config) LoadFromFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}
	return c.LoadFromContent(string(data))
}

// LoadFromContent 从内容加载配置.
func (c *Config) LoadFromContent(content string) error {
	return c.parseConfig(content)
}

// parseConfig 解析配置内容.
func (c *Config) parseConfig(content string) error {
	content = normalizeConfig(content)

	var section string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 检测段标记
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.ToLower(strings.Trim(line, "[]"))
			continue
		}

		// 解析键值对
		key, value, ok := parseKeyValue(line)
		if !ok {
			continue
		}

		switch section {
		case "interface":
			c.parseInterfaceField(key, value)
		case "peer":
			c.parsePeerField(key, value)
		}
	}

	c.Init()
	return nil
}

// parseInterfaceField 解析 Interface 段字段
func (c *Config) parseInterfaceField(key, value string) {
	switch strings.ToLower(key) {
	case "privatekey":
		c.PrivateKey = value
	case "address":
		for _, addr := range splitValues(value) {
			c.LocalAddrs = append(c.LocalAddrs, addr)
		}
	case "dns":
		c.DNSAddrs = splitValues(value)
	}
}

// parsePeerField 解析 Peer 段字段
func (c *Config) parsePeerField(key, value string) {
	switch strings.ToLower(key) {
	case "publickey":
		c.PublicKey = value
	case "presharedkey":
		c.PresharedKey = value
	case "endpoint":
		c.Endpoint = value
	case "allowedips":
		c.AllowedIPs = splitValues(value)
	case "persistentkeepalive":
		if v, err := strconv.Atoi(value); err == nil {
			c.PersistentKeepalive = v
		}
	}
}

// normalizeConfig 标准化配置格式，支持单行配置.
func normalizeConfig(content string) string {
	content = strings.ReplaceAll(content, "[Interface]", "\n[Interface]\n")
	content = strings.ReplaceAll(content, "[Peer]", "\n[Peer]\n")
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "[") {
			result = append(result, line)
			continue
		}
		if strings.Count(line, "=") > 1 {
			parts := splitMultipleKV(line)
			result = append(result, parts...)
		} else {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// splitMultipleKV 分割同一行中的多个键值对.
func splitMultipleKV(line string) []string {
	keys := []string{"PrivateKey", "PublicKey", "PresharedKey", "Address", "DNS",
		"Endpoint", "AllowedIPs", "PersistentKeepalive", "ListenPort"}
	var result []string
	remaining := line

	for len(remaining) > 0 {
		// 找到最早出现的 key
		minIdx := -1
		matchedKey := ""
		for _, key := range keys {
			idx := strings.Index(remaining, key)
			if idx != -1 && (minIdx == -1 || idx < minIdx) {
				minIdx = idx
				matchedKey = key
			}
		}
		if minIdx == -1 {
			break
		}
		// 跳过当前 key，找下一个 key 的位置
		afterKey := remaining[minIdx+len(matchedKey):]
		nextKeyIdx := -1
		for _, key := range keys {
			idx := strings.Index(afterKey, key)
			if idx != -1 && (nextKeyIdx == -1 || idx < nextKeyIdx) {
				nextKeyIdx = idx
			}
		}
		if nextKeyIdx == -1 {
			// 没有下一个 key，剩余部分就是当前键值对
			result = append(result, strings.TrimSpace(remaining[minIdx:]))
			break
		}
		// 提取当前键值对
		kv := remaining[minIdx : minIdx+len(matchedKey)+nextKeyIdx]
		result = append(result, strings.TrimSpace(kv))
		remaining = remaining[minIdx+len(matchedKey)+nextKeyIdx:]
	}
	return result
}

// parseKeyValue 解析键值对
func parseKeyValue(line string) (key, value string, ok bool) {
	key, value, ok = strings.Cut(line, "=")
	if !ok {
		return "", "", false
	}
	return strings.TrimSpace(key), strings.TrimSpace(value), true
}

// splitValues 分割逗号分隔的值
func splitValues(value string) []string {
	var result []string
	for _, v := range strings.Split(value, ",") {
		if v = strings.TrimSpace(v); v != "" {
			result = append(result, v)
		}
	}
	return result
}

// GetLocalNetipAddrs 获取本地地址
func (c *Config) GetLocalNetipAddrs() []netip.Addr {
	var addrs []netip.Addr
	for _, addr := range c.LocalAddrs {
		// 移除 CIDR 后缀
		if idx := strings.Index(addr, "/"); idx != -1 {
			addr = addr[:idx]
		}
		if ip, err := netip.ParseAddr(addr); err == nil {
			addrs = append(addrs, ip)
		}
	}
	return addrs
}

// GetDNSNetipAddrs 获取 DNS 地址.
func (c *Config) GetDNSNetipAddrs() []netip.Addr {
	var addrs []netip.Addr
	for _, addr := range c.DNSAddrs {
		if ip, err := netip.ParseAddr(addr); err == nil {
			addrs = append(addrs, ip)
		}
	}
	return addrs
}

// ToIPC 转换为 WireGuard IPC 格式
func (c *Config) ToIPC() string {
	lines := make([]string, 0)

	// Interface 部分
	if c.PrivateKey != "" {
		lines = append(lines, "private_key="+kcrypto.EncodeKeyHex(c.PrivateKey))
	}

	// Peer 部分
	if c.PublicKey != "" {
		lines = append(lines, "public_key="+kcrypto.EncodeKeyHex(c.PublicKey))
	}
	if c.PresharedKey != "" {
		lines = append(lines, "preshared_key="+kcrypto.EncodeKeyHex(c.PresharedKey))
	}
	if c.Endpoint != "" {
		lines = append(lines, "endpoint="+kurl.ResolveEndpoint(c.Endpoint))
	}
	for _, ip := range c.AllowedIPs {
		lines = append(lines, "allowed_ip="+ip)
	}
	if c.PersistentKeepalive > 0 {
		lines = append(lines, fmt.Sprintf("persistent_keepalive_interval=%d", c.PersistentKeepalive))
	}

	return strings.Join(lines, "\n")
}
