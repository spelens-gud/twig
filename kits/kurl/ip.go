package kurl

import "net"

// ResolveEndpoint 解析端点地址，将域名解析为 IP.
func ResolveEndpoint(endpoint string) string {
	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		return endpoint
	}
	if ip := net.ParseIP(host); ip != nil {
		return endpoint
	}
	// 解析域名
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		return endpoint
	}
	// 使用ipv4
	for _, ip := range ips {
		if ip.To4() != nil {
			return net.JoinHostPort(ip.String(), port)
		}
	}
	return net.JoinHostPort(ips[0].String(), port)
}
