package main

import (
	"log"
	"shenlong/ippool_server/server"
)

//"shenlong/ippool_server/ippoolserver"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	appConfig := server.GetAppConf()
	appConfig.LoadConfig()
	server.Start()
}
