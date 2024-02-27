package handler

import (
	"dstat/pkg/ws"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Allenxuxu/gev"
)

type DStat struct {
	ConnPerSec atomic.Int64
	ActiveConn atomic.Int64
	Inbound    atomic.Int64

	IPsMap    map[string]bool
	map_mutex sync.Mutex

	wss *ws.WebSocketServer
}

func NewDStat(wss *ws.WebSocketServer) *DStat {
	return &DStat{
		IPsMap: make(map[string]bool),
		wss:    wss,
	}
}

func (d *DStat) OnClose(c *gev.Connection) {
	d.ActiveConn.Add(-1)
}

func (d *DStat) OnConnect(c *gev.Connection) {
	d.ConnPerSec.Add(1)
	d.ActiveConn.Add(1)

	d.map_mutex.Lock()
	d.IPsMap[strings.Split(c.PeerAddr(), ":")[0]] = true
	d.map_mutex.Unlock()
}

func (d *DStat) OnMessage(c *gev.Connection, ctx interface{}, data []byte) (out interface{}) {
	d.Inbound.Add(int64(len(data)))
	return
}

func (d *DStat) BroadcastAndReset() {
	stats := fmt.Sprintf(`{"numConnPerSec": %d, "numActiveConn": %d, "inboundMBps": %f, "numIpsPerSec": %d}`,
		d.ConnPerSec.Swap(0),
		d.ActiveConn.Load(),
		(float32(d.Inbound.Swap(0)) / 1024.0 / 1024.0),
		len(d.IPsMap))

	d.wss.Broadcast([]byte(stats))

	d.map_mutex.Lock()
	d.IPsMap = make(map[string]bool)
	d.map_mutex.Unlock()
}
