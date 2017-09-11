package main

import (
	"12306/httprequest"
	"12306/log"
	"flag"
)

var (
	usename  *string = flag.String("username", "", "login usename")
	password *string = flag.String("password", "", "login password")
)

func main() {
	flag.Parse()

	client := httprequest.NewClient()
	err := client.Login(*usename, *password)
	if err != nil {
		return
	}
	ps, _ := client.GetPassengers()
	log.MyLogInfo("passenger:%v\n", ps)
}
