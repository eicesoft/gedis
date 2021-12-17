package main

import (
	"flag"
	"fmt"
	"gedis/config"
	"gedis/pkg/logger"
	"gedis/server"
)

var (
	host string
	port string
)

func init() {
	flag.StringVar(&host, "host", "", "server listen host")
	flag.StringVar(&port, "port", config.Get().Server.Port, "server listen port")
	flag.Parse()
}

func main() {
	addr := fmt.Sprintf("%s:%s", host, port)
	s := server.NewServer(addr)
	logger.Info("Gedis listen in: " + addr)
	s.Run()
}
