package iproxy

import (
	"context"
	"net"
	"net/netip"
	"os"

	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
)

// NetTun 网络栈接口.
type NetTun interface {
	Name() (string, error)
	File() *os.File
	Read([]byte, int) (int, error)
	Write([]byte, int) (int, error)
	Flush() error
	Close() error
	MTU() (int, error)
}

// Net 抽象网络接口.
type Net interface {
	ListenTCP(addr *net.TCPAddr) (*gonet.TCPListener, error)
	DialUDPAddrPort(laddr, raddr netip.AddrPort) (*gonet.UDPConn, error)
	ListenUDPAddrPort(laddr netip.AddrPort) (*gonet.UDPConn, error)
	DialUDP(laddr, raddr *net.UDPAddr) (*gonet.UDPConn, error)
	ListenUDP(laddr *net.UDPAddr) (*gonet.UDPConn, error)
	LookupHost(host string) (addrs []string, err error)
	LookupContextHost(ctx context.Context, host string) ([]string, error)
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
	Dial(network, address string) (net.Conn, error)
	DialContextTCP(ctx context.Context, addr *net.TCPAddr) (*gonet.TCPConn, error)
}
