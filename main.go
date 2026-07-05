package main

import (
	"flag"
	"log"
	"github.com/xd-sarthak/go-redis/config"
	"github.com/xd-sarthak/go-redis/server"
)

func setUpFlags(){
	flag.StringVar(&config.Host,"host","0.0.0.0","host for the server ")
	flag.IntVar(&config.Port,"port",7379,"port for the server")
	flag.Parse()
}

func main() {
	setUpFlags()
	log.Println("starting the server...")
	server.RunAsyncTCPServer()
}