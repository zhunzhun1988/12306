package main

import (
	"12306/httprequest"
	//"12306/log"
)

func main() {
	client := httprequest.NewClient()
	err := client.Login("***@qq.com", "####")
	if err != nil {
		return
	}
	client.GetPassengers()
}
