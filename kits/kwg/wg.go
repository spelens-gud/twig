// Package kwg 提供 WireGuard VPN 代理功能
// 基于用户态网络栈实现，无需 root 权限即可运行
package kwg

import (
	"context"
	"fmt"

	"github.com/spelens-gud/twig/interfaces/iproxy"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

var (
	wgNet     *netstack.Net
	wgDevice  iproxy.Device
	running   bool
	allowedIP []string
)

// Run 快捷启动方法，创建并启动代理
func Run(ctx context.Context, path string) error {
	cfg := NewConfig()
	err := cfg.LoadFromFile(path)
	if err != nil {
		// TODO 日志模块
		return fmt.Errorf("加载配置文件失败: %w", err)
	}
	cfg.Init()

	if running {
		return fmt.Errorf("代理已在运行中")
	}

	allowedIP = cfg.AllowedIPs

	// 创建用户态 TUN 接口
	tunDev, gNet, err := netstack.CreateNetTUN(cfg.GetLocalNetipAddrs(), cfg.GetDNSNetipAddrs(), device.DefaultMTU)
	if err != nil {
		return fmt.Errorf("创建 TUN 接口失败: %w", err)
	}

	// 创建 WireGuard 设备
	logger := device.NewLogger(device.LogLevelError, "wireguard: ")
	wgDev := device.NewDevice(tunDev, conn.NewDefaultBind(), logger)

	// 将配置转换为 IPC 格式并应用
	ipcConf := cfg.ToIPC()
	if err := wgDev.IpcSet(ipcConf); err != nil {
		return fmt.Errorf("应用 WireGuard 配置失败: %w", err)
	}

	// 启动设备
	if err := wgDev.Up(); err != nil {
		return fmt.Errorf("启动 WireGuard 设备失败: %w", err)
	}

	wgDevice = wgDev
	wgNet = gNet
	running = true

	// 监听上下文取消
	go func() {
		<-ctx.Done()
		//nolint:errcheck
		Stop()
	}()

	return nil
}

// Stop 停止 WireGuard 代理.
func Stop() error {
	if !running {
		return nil
	}

	if wgDevice != nil {
		wgDevice.Close()
		wgDevice = nil
	}

	wgNet = nil
	running = false
	return nil
}

// IsRunning 检查代理是否正在运行.
func IsRunning() bool {
	return running
}

// Net 获取网络栈.
func Net() *netstack.Net {
	return wgNet
}

// AllowedIPs 获取允许的 IP 范围.
func AllowedIPs() []string {
	return allowedIP
}
