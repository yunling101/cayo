package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yunling101/cayo/pkg/cache"
	"github.com/yunling101/cayo/pkg/cayoserver"
	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/model/user"
)

func main() {
	version := flag.Bool("v", false, "show version")
	cert := flag.String("cert", "certs/server.pem", "path to PEM encoded public key certificate.")
	key := flag.String("key", "certs/server.key", "path to private key associated with given certificate.")
	cfg := flag.String("c", "config/server.yml", "config file")
	flag.Parse()

	if *version {
		fmt.Println(global.Version)
		os.Exit(0)
	}
	log.Println("starting cayoserver version", global.Version)

	global.LoadConfig(*cfg)
	global.NewDB()
	cache.NewInit()

	go cayoserver.NewRpcStart(*cert, *key)
	go cayoserver.NewHTTPStart()
	go user.InitUserAdmin()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sc:
		global.Session.Close()
		os.Exit(0)
	}
}
