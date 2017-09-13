package httprequest

import (
	"12306/utils"
	"12306/verifycode"
	"fmt"
	"net/url"
	"sort"
)

const (
	host_addr                   = "kyfw.12306.cn"
	login_init_addr             = string("https://") + string(host_addr) + string("/otn/login/init")
	login_init12306_addr        = string("https://") + string(host_addr) + string("/otn/index/initMy12306")
	login_verify_codeimage_addr = string("https://") + string(host_addr) + string("/passport/captcha/captcha-image?login_site=E&module=login&rand=sjrand&")
	login_verify_addr           = string("https://") + string(host_addr) + string("/passport/captcha/captcha-check")
	weblogin_addr               = string("https://") + string(host_addr) + string("/passport/web/login")
	userlogin_addr1             = string("https://") + string(host_addr) + string("/otn/login/userLogin")
	userlogin_addr2             = string("https://") + string(host_addr) + string("/otn/passport?redirect=/otn/login/userLogin")
	userlogin_check             = string("https://") + string(host_addr) + string("/otn/login/checkUser")

	get_token_addr = string("https://") + string(host_addr) + string("/passport/web/auth/uamtk")
	set_token_addr = string("https://") + string(host_addr) + string("/otn/uamauthclient")

	get_station_addr            = string("https://") + string(host_addr) + string("/otn/resources/js/framework/station_name.js?station_version=1.9025")
	leftticket_init_addr        = string("https://") + string(host_addr) + string("/otn/leftTicket/init")
	get_leftticket_addr         = string("https://") + string(host_addr) + string("/otn/leftTicket/queryX")
	leftticket_logindevice_addr = string("https://") + string(host_addr) + string("/otn/HttpZF/logdevice?algID=P2z9BuCc7X&hashCode=SlGo_d6sfp9xtw8AtZ_duN98eWcL3DKpPKIJmoBWPu0&FMQw=0&q4f3=zh-CN&VySQ=FFEVRmrfjtNaxt7gStfKnV1e-COX_t4t&VPIf=1&custID=133&VEek=unknown&dzuS=0&yD16=0&EOQP=4902a61a235fbb59700072139347967d&lEnu=169052788&jp76=f02d9c91345cd461956c69d8807f1b23&hAqN=Win32&platform=WEB&ks0Q=e6917e2a69332dc7f73f1e97d89f42d8&TeRS=1023x2037&tOHY=24xx1063x2037&Fvje=i1l1o1s1&q5aJ=-8&wNLf=99115dfb07133750ba677d055874de87&0aew=Mozilla/5.0%20(Windows%20NT%2010.0;%20WOW64)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/60.0.3112.113%20Safari/537.36&E3gR=d69c1ff1a8305f9ca801372d1c87366d")
	leftticket_log_addr         = string("https://") + string(host_addr) + string("/otn/leftTicket/log")
)

func getLeftTicketLoginDeviceUrl() string {
	return fmt.Sprintf("%s&timestamp=%s", leftticket_logindevice_addr, utils.GetNowMicoSecondStr())
}

func getLeftTicketUrl(date, fromStation, toStation, code string) string {
	return fmt.Sprintf("%s?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=%s",
		get_leftticket_addr, date, fromStation, toStation, code)
}

func getLeftTicketLogUrl(date, fromStation, toStation, code string) string {
	return fmt.Sprintf("%s?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=%s",
		leftticket_log_addr, date, fromStation, toStation, code)
}

func getLoginVerifyImgUrl() string {
	return string(login_verify_codeimage_addr) + utils.GetRandFloat(16)
}

func getLoginVerifyUrlValues(poss verifycode.VerifyPosList) url.Values {
	sort.Sort(poss)
	ret := make(url.Values)
	ret["login_site"] = []string{"E"}
	ret["rand"] = []string{"sjrand"}
	ret["answer"] = []string{poss.ToString()}
	return ret
}

func getLoginUrlValues(usrname, password string) url.Values {
	ret := make(url.Values)
	ret["username"] = []string{usrname}
	ret["password"] = []string{password}
	ret["appid"] = []string{"otn"}
	return ret
}
