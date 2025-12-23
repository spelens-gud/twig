package kwg

import (
	"context"
	"fmt"
	"net"
)

// Dialer 通过 WireGuard 隧道建立连接的拨号器
type Dialer struct {
}

// NewDialer 创建新的拨号器
func NewDialer() *Dialer {
	return &Dialer{}
}

// DialContext 通过 WireGuard 隧道建立 TCP 连接
func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if !IsRunning() {
		return nil, fmt.Errorf("WireGuard 代理未运行")
	}

	gNet := Net()
	if gNet == nil {
		return nil, fmt.Errorf("网络栈不可用")
	}

	switch network {
	case "tcp", "tcp4", "tcp6":
		return gNet.DialContextTCP(ctx, parseTCPAddr(address))
	case "udp", "udp4", "udp6":
		return gNet.DialUDP(nil, parseUDPAddr(address))
	default:
		return nil, fmt.Errorf("不支持的网络类型: %s", network)
	}
}

// Dial 通过 WireGuard 隧道建立连接
func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

// parseTCPAddr 解析 TCP 地址
func parseTCPAddr(address string) *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil
	}
	return addr
}

// parseUDPAddr 解析 UDP 地址
func parseUDPAddr(address string) *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil
	}
	return addr
}
