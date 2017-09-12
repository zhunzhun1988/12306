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
	get_token_addr              = string("https://") + string(host_addr) + string("/passport/web/auth/uamtk")
	set_token_addr              = string("https://") + string(host_addr) + string("/otn/uamauthclient")

	get_station_addr    = string("https://") + string(host_addr) + string("/otn/resources/js/framework/station_name.js?station_version=1.9025")
	set_leftticket_addr = string("https://") + string(host_addr) + string("/otn/leftTicket/queryX")
)

func getLeftTicketUrl(date, fromStation, toStation, code string) string {
	return fmt.Sprintf("%s?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=%s",
		set_leftticket_addr, date, fromStation, toStation, code)
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
