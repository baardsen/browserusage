package main

import (
	"browserusage/webserver"
	"browserusage/dao"
	"flag"
)

var (
	port *int  = flag.Int("port", 9050, "Port number")
)

func main() {
	dao.Init()
	flag.Parse()
	webserver.Start(*port)
}