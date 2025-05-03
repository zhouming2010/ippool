package main

import (
	"log"
	"os"
	"shenlong/ippool_server/proxyserver"
)

func main() {
	log.SetOutput(os.Stdout)

	proxyserver.GetInstance().Start()

}
