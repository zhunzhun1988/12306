package main

import (
	"12306/httprequest"
	//"12306/utils"
	"flag"
	"fmt"
	"time"
)

var (
	usename  *string = flag.String("username", "", "login usename")
	password *string = flag.String("password", "", "login password")
)

func main() {
	flag.Parse()
	client := httprequest.NewClient(*usename, *password)

	_, wait, err := client.CheckAndOrderTicket("2017-10-12", "上海", "嘉兴", []string{"谢谆志"}, []string{"K1805"}, []httprequest.TicketType{httprequest.Ticket_YZ}, time.Second)
	if err != nil {
		return
	}
	if wait != nil {
		ok, err := wait(100 * time.Second)
		fmt.Printf("exit, %t, %v\n", ok, err)
	}
}
