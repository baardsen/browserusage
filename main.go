package main

import (
	"browserusage/dao"
	"browserusage/webserver"
	"flag"
)

var (
	port *int = flag.Int("port", 9050, "Port number")
)

func main() {
	dao.Init()
	flag.Parse()
	webserver.Start(*port)
}
