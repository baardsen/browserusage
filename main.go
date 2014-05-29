package main

import (
	"browserusage/dao"
	resourcelocator "github.com/baardsen/resourcelocator"
	"browserusage/webserver"
	"flag"
	"fmt"
	"time"
)

var (
	port *int  = flag.Int("port", 9050, "Port number")
)

func main() {
	flag.Parse() // Scan the arguments list
	if flag.NArg() == 0 {
		webserver.Start(*port)
	} else {
		from := time.Date(2014, 04, 7, 0, 0, 0, 0, time.Local)
		fmt.Printf("%+v", dao.Query(from, time.Now()))
		//flag.Usage()
	}
}