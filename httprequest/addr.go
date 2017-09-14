package httprequest

import (
	"12306/utils"
	"12306/verifycode"
	"fmt"
	"net/url"
	"sort"
	"strings"
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

	order_ticket_addr        = string("https://") + string(host_addr) + string("/otn/leftTicket/submitOrderRequest")
	get_submittoken_addr     = string("https://") + string(host_addr) + string("/otn/confirmPassenger/initDc")
	check_order_addr         = string("https://") + string(host_addr) + string("/otn/confirmPassenger/checkOrderInfo")
	get_orderqueuecount_addr = string("https://") + string(host_addr) + string("/otn/confirmPassenger/getQueueCount")
	confirm_order_addr       = string("https://") + string(host_addr) + string("/otn/confirmPassenger/confirmSingleForQueue")
)

func getOrderTickerUrlValueStr(secret, date, backdate, from, to, flag, code string) string {
	return fmt.Sprintf("secretStr=%s&train_date=%s&back_train_date=%s&tour_flag=%s&purpose_codes=%s&query_from_station_name=%s&query_to_station_name=%s&undefined",
		secret, date, backdate, flag, code, from, to)
}

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
func getConfirmOrderUrlValuesStr(keycheck, leftticket, token string, ps []Passenger, st SeatType, pt PassengerType) string {
	strPTS := make([]string, 0, len(ps))
	strOPTS := make([]string, 0, len(ps))
	for _, passenger := range ps {
		strPTS = append(strPTS, fmt.Sprintf("%s,0,%d,%s,%s,%s,%s,N", st, pt, passenger.PassengerName, passenger.PassengerIDTypeCode, passenger.PassengerIDNo, passenger.MobileNo))
		strOPTS = append(strOPTS, fmt.Sprintf("%s,%s,%s,%d", passenger.PassengerName, passenger.PassengerIDTypeCode, passenger.PassengerIDNo, pt))
	}

	return fmt.Sprintf("passengerTicketStr=%s&oldPassengerStr=%s_&randCode=&purpose_codes=00&key_check_isChange=%s&leftTicketStr=%s&train_location=J1&choose_seats=&seatDetailType=000&roomType=00&dwAll=N&_json_att=&REPEAT_SUBMIT_TOKEN=%s",
		utils.UrlEncode(strings.Join(strPTS, "_")), utils.UrlEncode(strings.Join(strOPTS, "_")), keycheck, utils.UrlEncode(leftticket), token)

	//return fmt.Sprintf("randCode=&purpose_codes=00&key_check_isChange=%s&leftTicketStr=%s&train_location=J1&choose_seats=&seatDetailType=000&roomType=00&dwAll=N&_json_att=&REPEAT_SUBMIT_TOKEN=%s&passengerTicketStr=%s&oldPassengerStr=%s",

	//keycheck, utils.UrlEncode(leftticket), token, utils.UrlEncode(strings.Join(strPTS, "_")), utils.UrlEncode(strings.Join(strOPTS, "_")))
}

func getConfirmOrderUrlValues(keycheck, leftticket, token string, ps []Passenger, st SeatType, pt PassengerType) url.Values {
	ret := make(url.Values)
	ret["randCode"] = []string{""}
	ret["purpose_codes"] = []string{"00"}
	ret["key_check_isChange"] = []string{keycheck}
	ret["leftTicketStr"] = []string{leftticket}
	ret["train_location"] = []string{"J1"}
	ret["choose_seats"] = []string{""}
	ret["seatDetailType"] = []string{"000"}
	ret["roomType"] = []string{"00"}
	ret["dwAll"] = []string{"N"}
	ret["_json_att"] = []string{""}
	ret["REPEAT_SUBMIT_TOKEN"] = []string{token}
	strPTS := make([]string, 0, len(ps))
	strOPTS := make([]string, 0, len(ps))
	for _, passenger := range ps {
		strPTS = append(strPTS, fmt.Sprintf("%s,0,%d,%s,%s,%s,%s,N", st, pt, passenger.PassengerName, passenger.PassengerIDTypeCode, passenger.PassengerIDNo, passenger.MobileNo))
		strOPTS = append(strOPTS, fmt.Sprintf("%s,%s,%s,%d", passenger.PassengerName, passenger.PassengerIDTypeCode, passenger.PassengerIDNo, pt))
	}
	ret["passengerTicketStr"] = []string{strings.Join(strPTS, "_")}
	ret["oldPassengerStr"] = []string{strings.Join(strOPTS, "_") + "_"}
	return ret
}

func getCheckOrderUrlValues(cancelflag, orderNum, flag, token string, ps []Passenger, st SeatType, pt PassengerType) url.Values {
	ret := make(url.Values)
	ret["cancel_flag"] = []string{cancelflag}
	ret["bed_level_order_num"] = []string{orderNum}
	ret["tour_flag"] = []string{flag}
	ret["randCode"] = []string{""}
	ret["_json_att"] = []string{""}
	ret["REPEAT_SUBMIT_TOKEN"] = []string{token}
	strPTS := make([]string, 0, len(ps))
	strOPTS := make([]string, 0, len(ps))
	for _, passenger := range ps {
		strPTS = append(strPTS, fmt.Sprintf("%s,0,%d,%s,%s,%s,%s,N", st, pt, passenger.PassengerName, passenger.PassengerIDTypeCode, passenger.PassengerIDNo, passenger.MobileNo))
		strOPTS = append(strOPTS, fmt.Sprintf("%s,%s,%s,%d", passenger.PassengerName, passenger.PassengerIDTypeCode, passenger.PassengerIDNo, pt))
	}
	ret["passengerTicketStr"] = []string{strings.Join(strPTS, "_")}
	ret["oldPassengerStr"] = []string{strings.Join(strOPTS, "_")}
	return ret
}

func getOrderQueueCountUrlValues(date, trainno, trainCode, setType, fromStationCode, toStationCode, leftTicket, token string) url.Values {
	ret := make(url.Values)
	ret["train_date"] = []string{date}
	ret["train_no"] = []string{trainno}
	ret["stationTrainCode"] = []string{trainCode}
	ret["seatType"] = []string{setType}
	ret["fromStationTelecode"] = []string{fromStationCode}
	ret["toStationTelecode"] = []string{toStationCode}
	ret["leftTicket"] = []string{leftTicket}
	ret["purpose_codes"] = []string{"00"}
	ret["train_location"] = []string{"J1"}
	ret["_json_att"] = []string{""}
	ret["REPEAT_SUBMIT_TOKEN"] = []string{token}
	return ret
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
