package main

import (
	"browserusage/webserver"
	"flag"
)

var (
	port *int  = flag.Int("port", 9050, "Port number")
)

func main() {
	flag.Parse()
	webserver.Start(*port)
}