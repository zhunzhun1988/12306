package main

import (
	"12306/httprequest"
	//"encoding/json"
	"time"
	//"time"
	//"12306/log"
	//"12306/utils"
	"flag"
	"fmt"
)

var (
	usename  *string = flag.String("username", "", "login usename")
	password *string = flag.String("password", "", "login password")
)

func main() {
	flag.Parse()
	client := httprequest.NewClient(*usename, *password)

	cancel := client.CheckAndOrderTicket("2017-10-12", "上海", "嘉兴", []string{"K1805", "G7301"}, httprequest.Ticket_YW, time.Second)
	time.Sleep(20 * time.Second)
	cancel()
	time.Sleep(10 * time.Second)
	fmt.Printf("exit")

}
