package main

import (
	"dstat/pkg/handler"
	"dstat/pkg/ws"
	"flag"
	"net/http"
	"time"

	"github.com/Allenxuxu/gev"
)

var (
	loops         = flag.Int("l", 1, "loops")
	addr          = flag.String("addr", "0.0.0.0:1234", "server address")
	webSocketAddr = flag.String("wsa", "0.0.0.0:8080", "WebSocket server address")
)

func main() {
	flag.Parse()

	wss := ws.NewWebSocketServer(*webSocketAddr)
	wss.Start()

	d := handler.NewDStat(wss)

	s, err := gev.NewServer(d,
		gev.Network("tcp"),
		gev.Address(*addr),
		gev.NumLoops(*loops),
		gev.IdleTime(time.Second*5),
	)

	if err != nil {
		panic(err)
	}

	s.RunEvery(time.Second, d.BroadcastAndReset)

	go HTTPServer()

	s.Start()
}

func HTTPServer() {
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.ListenAndServe(":80", nil)
}
