package iproxy

import (
	"golang.zx2c4.com/wireguard/conn"
)

type Device interface {
	Down() error
	IsUnderLoad() bool
	RemoveAllPeers()
	Close()
	Wait() chan struct{}
	SendKeepalivesToPeersWithCurrentKeypair()
	Bind() conn.Bind
	BindSetMark(mark uint32) error
	BindUpdate() error
	BindClose() error
}
