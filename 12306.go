package main

import (
	"12306/httprequest"
	"strings"
	//"12306/utils"
	"12306/log"
	"flag"
	"fmt"
	"time"
)

var (
	usename     *string = flag.String("username", "", "12306登录用户名")
	password    *string = flag.String("password", "", "12306登录密码")
	orderDate   *string = flag.String("date", "", "要购买车票日期:如 2017-10-01")
	formStation *string = flag.String("fromstation", "", "出发车站：如　上海")
	toStation   *string = flag.String("tostation", "", "目的车站：　如　北京")
	trainNo     *string = flag.String("trainnum", "", "车次：　如　G102")
	passenger   *string = flag.String("passenger", "", "为谁购票：　如　张三")
	seatType    *string = flag.String("seattype", "", "车票座位类型：　如　硬卧")
	timeout     *int    = flag.Int("timeout", 1000000, "监控余票超时时间（秒）")
	interval    *int    = flag.Int("interval", 1000, "扫描频率（毫秒）")
	debug       *bool   = flag.Bool("debug", false, "是否开启调试模式")
)

func checkArgs() error {
	if *usename == "" || *password == "" {
		return fmt.Errorf("请输入12306登录用户名和密码")
	}
	if *orderDate == "" {
		return fmt.Errorf("请输入要购买车票的日期：　如　2017-10-01")
	} else {
		p := strings.Split(*orderDate, "-")
		if len(p) != 3 {
			return fmt.Errorf("请输入正确购票日期：　如　2017-10-01")
		}

	}
	if *formStation == "" || *toStation == "" {
		return fmt.Errorf("请输入出发站和目的站")
	}
	if *trainNo == "" {
		return fmt.Errorf("请输入购买车次")
	}
	if *passenger == "" {
		return fmt.Errorf("请输入乘客名字")
	}
	if *seatType == "" || httprequest.StringToSeatType(*seatType) == httprequest.Ticket_UNKNOW {
		return fmt.Errorf("请输入正确座位类型")
	}
	return nil
}
func main() {
	flag.Parse()

	if err := checkArgs(); err != nil {
		log.MyLogE("%s", err.Error())
		return
	}
	log.SetDebug(*debug)
	client := httprequest.NewClient(*usename, *password)

	_, wait, err := client.CheckAndOrderTicket(*orderDate, *formStation, *toStation, []string{*passenger}, []string{*trainNo}, []httprequest.TicketType{httprequest.StringToSeatType(*seatType)}, time.Millisecond*time.Duration(*interval))
	if err != nil {
		return
	}
	if wait != nil {
		_, err := wait(time.Duration(*timeout) * time.Second)
		log.MyLogI("%s", err)
	}
}
