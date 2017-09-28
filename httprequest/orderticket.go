package httprequest

import (
	"12306/log"
	"12306/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TicketType int

const (
	Ticket_UNKNOW = -1
	Ticket_TDZ    = iota
	Ticket_YDZ
	Ticket_EDZ
	Ticket_GJRW
	Ticket_RW
	Ticket_DW
	Ticket_YW
	Ticket_RZ
	Ticket_YZ
	Ticket_WZ
)

type PassengerType int

const (
	ADULT = 1
	CHILDREN
	STUDENT
	SCJR
)

type SeatType string

const (
	SeatType_UNKNOW = ""
	SeatType_YZ     = "1"
	SeatType_WZ     = "1"
	SeatType_YW     = "3"
	SeatType_RW     = "4"
	SeatType_EDZ    = "O"
	SeatType_YDZ    = "M"
	SeatType_SWZ    = "9"
	SeatType_GJRW   = "6"
)

func StringToSeatType(str string) TicketType {
	switch str {
	case "特等座":
		return Ticket_TDZ
	case "一等座":
		return Ticket_YDZ
	case "二等座":
		return Ticket_EDZ
	case "高级软卧":
		return Ticket_GJRW
	case "软卧":
		return Ticket_RW
	case "动卧":
		return Ticket_DW
	case "硬卧":
		return Ticket_YW
	case "软座":
		return Ticket_RZ
	case "硬座":
		return Ticket_YZ
	case "无座":
		return Ticket_WZ
	}
	return Ticket_UNKNOW
}

func ticketTypeTSeatType(tt TicketType) SeatType {
	switch tt {
	case Ticket_YZ:
		return SeatType_YZ
	case Ticket_WZ:
		return SeatType_WZ
	case Ticket_YW:
		return SeatType_YW
	case Ticket_RW:
		return SeatType_RW
	case Ticket_EDZ:
		return SeatType_EDZ
	case Ticket_YDZ:
		return SeatType_YDZ
	case Ticket_GJRW:
		return SeatType_GJRW
	}
	return SeatType_UNKNOW
}
func isTicketMatchType(ti *TicketsInfo, types []TicketType) (bool, TicketType) {
	if ti.HaveTickets == false {
		return false, Ticket_WZ
	}
	for _, t := range types {
		switch t {
		case Ticket_TDZ:
			if ti.TDZ != "" && ti.TDZ != "无" {
				return true, Ticket_TDZ
			}
		case Ticket_YDZ:
			if ti.YDZ != "" && ti.YDZ != "无" {
				return true, Ticket_YDZ
			}
		case Ticket_EDZ:
			if ti.EDZ != "" && ti.EDZ != "无" {
				return true, Ticket_EDZ
			}
		case Ticket_GJRW:
			if ti.GJRW != "" && ti.GJRW != "无" {
				return true, Ticket_GJRW
			}
		case Ticket_RW:
			if ti.RW != "" && ti.RW != "无" {
				return true, Ticket_RW
			}
		case Ticket_DW:
			if ti.DW != "" && ti.DW != "无" {
				return true, Ticket_DW
			}
		case Ticket_YW:
			if ti.YW != "" && ti.YW != "无" {
				return true, Ticket_YW
			}
		case Ticket_RZ:
			if ti.RZ != "" && ti.RZ != "无" {
				return true, Ticket_RZ
			}
		case Ticket_YZ:
			if ti.YZ != "" && ti.YZ != "无" {
				return true, Ticket_YZ
			}
		case Ticket_WZ:
			if ti.WZ != "" && ti.WZ != "无" {
				return true, Ticket_WZ
			}
		}
	}
	return false, Ticket_WZ
}
func OrderTicket(client *http.Client, secret, date, from, to string) error {
	resp, err := client.Post(order_ticket_addr, "application/x-www-form-urlencoded; charset=UTF-8",
		strings.NewReader(getOrderTickerUrlValueStr(secret, date, date, from, to, "dc", "ADULT")))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OrderTicket bad status code:%d", resp.StatusCode)
	}
	body := getBody(resp.Body)
	log.MyOrderLogD("OrderTicket:body[%s]", string(body))
	if len(body) == 0 {
		return fmt.Errorf("OrderTicket status is empty")
	}
	otm := OrderTicketMsg{}
	err = json.Unmarshal(body, &otm)
	if err != nil {
		return fmt.Errorf("OrderTicket json Unmarshal err:%v,[%s]", err, string(body))
	}
	if otm.Status == false {
		return fmt.Errorf("OrderTicket fail:%s", strings.Join(otm.Messages, ","))
	}
	return nil
}

var myExp = regexp.MustCompile(`'key_check_isChange':'(?P<token>\w+)'`)
var myExp2 = regexp.MustCompile(`globalRepeatSubmitTokens\s+=\s+'(?P<token>\w+)';`)

func getGlobalRepeatSubmitTokenStr(buf string) (string, string) {
	checkToken := ""
	submitToken := ""
	strCheckTokens := myExp.FindStringSubmatch(buf)

	if len(strCheckTokens) == 2 {
		checkToken = strCheckTokens[1]
	}

	strs := strings.Split(buf, "\n")
	key := "globalRepeatSubmitToken"
	for _, str := range strs {
		if index := strings.Index(str, key); index > 0 {
			si := strings.Index(str, "'")
			ei := strings.LastIndex(str, "'")
			if si > 0 && ei > si {
				submitToken = str[si+1 : ei]
			}
		}
	}
	return checkToken, submitToken
}
func GetSubmitToken(client *http.Client) (string, string, error) {
	resp, err := client.PostForm(get_submittoken_addr, url.Values{"_json_att": []string{}})
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GetSubmitToken bad status code:%d", resp.StatusCode)
	}
	body := getBody(resp.Body)
	checktoken, submitToken := getGlobalRepeatSubmitTokenStr(string(body))
	if checktoken != "" && submitToken != "" {
		return checktoken, submitToken, nil
	}
	return "", "", fmt.Errorf("获取token失败")
}

func CheckOrderInfo(client *http.Client, passengers []Passenger, tt TicketType, token string) (ok, showPC bool, seatType SeatType, err error) {
	st := ticketTypeTSeatType(tt)
	if st == SeatType_UNKNOW {
		return false, false, SeatType_UNKNOW, fmt.Errorf("暂不支持该TicketType")
	}
	resp, err := client.PostForm(check_order_addr, getCheckOrderUrlValues("2", "000000000000000000000000000000", "dc", token, passengers, st, ADULT))
	if err != nil {
		return false, false, st, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, false, st, fmt.Errorf("CheckOrderInfo bad status code:%d", resp.StatusCode)
	}
	body := getBody(resp.Body)
	log.MyOrderLogD("CheckOrderInfo:body[%s]", string(body))
	com := CheckOrderMsg{}
	errJson := json.Unmarshal(body, &com)
	if errJson != nil {
		return false, false, st, fmt.Errorf("CheckOrderInfo json err:%v,[%s]", errJson, string(body))
	}
	if com.Status != true {
		return false, false, st, fmt.Errorf("CheckOrderInfo fail msg:%s", strings.Join(com.Messages, ","))
	}
	return true, com.Data.IfShowPassCode != "N", st, nil
}

func GetOrderQueueCount(client *http.Client, dateTime time.Time, trainno, trainCode, seatType, fromStationCode, toStationCode, leftTicket, token string) (ok bool, leftNum int, err error) {
	date := utils.GetOrderTimeFomat(dateTime, true)
	dataUrl := getOrderQueueCountUrlValues(date, trainno, trainCode, seatType, fromStationCode, toStationCode, leftTicket, token)
	resp, err := client.PostForm(get_orderqueuecount_addr, dataUrl)
	if err != nil {
		return false, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, 0, fmt.Errorf("CheckOrderInfo bad status code:%d", resp.StatusCode)
	}
	body := getBody(resp.Body)
	log.MyOrderLogD("GetOrderQueueCount:body[%s]", string(body))
	oqcm := OrderQueueCountMsg{}
	errJson := json.Unmarshal(body, &oqcm)
	if errJson != nil {
		return false, 0, fmt.Errorf("GetOrderQueueCount json err:%v,[%s]", errJson, string(body))
	}
	num := 0
	if seatType == SeatType_YZ || seatType == SeatType_WZ {
		strs := strings.Split(oqcm.Data.Ticket, ",")
		if len(strs) == 2 {
			if seatType == SeatType_YZ {
				num, _ = strconv.Atoi(strs[0])
			} else {
				num, _ = strconv.Atoi(strs[1])
			}
		}
	} else {
		num, _ = strconv.Atoi(oqcm.Data.Ticket)
	}
	if oqcm.Status == false {
		return false, 0, fmt.Errorf("GetOrderQueueCount fail msg:%s", strings.Join(oqcm.Messages, ","))
	}
	return true, num, nil
}

func ConfirmOrder(client *http.Client, keycheck, leftticket, token string, ps []Passenger, st SeatType) (ok bool, err error) {
	/*buf := utils.UrlEncode([]byte(getConfirmOrderUrlValuesStr(keycheck, leftticket, token, ps, st, ADULT)))
	fmt.Printf("buf:[%s]\n", buf)
	req, err := http.NewRequest("POST", confirm_order_addr, strings.NewReader(buf))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, errPost := client.Do(req)
	if errPost != nil {
		return false, errPost
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("ConfirmOrder bad status code:%d", resp.StatusCode)
	}
	body := getBody(resp.Body)
	log.MyOrderLogD("ConfirmOrder:body[%s]", string(body))*/

	req, err := http.NewRequest("POST", confirm_order_addr, strings.NewReader(getConfirmOrderUrlValuesStr(keycheck, leftticket, token, ps, st, ADULT)))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	resp, errPost := client.Do(req)
	//resp, errPost := client.Post(confirm_order_addr, "application/x-www-form-urlencoded; charset=UTF-8",
	//	strings.NewReader(getConfirmOrderUrlValuesStr(keycheck, leftticket, token, ps, st, ADULT)))
	if errPost != nil {
		return false, errPost
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("ConfirmOrder bad status code:%d", resp.StatusCode)
	}
	//body := getBody(resp.Body)
	//log.MyOrderLogD("ConfirmOrder:body[%s]", string(body))
	return true, nil

	/*dataUrl := getConfirmOrderUrlValues(keycheck, leftticket, token, ps, st, ADULT)
	fmt.Printf("urlDat:%v\n", dataUrl)
	resp, err := client.PostForm(confirm_order_addr, dataUrl)

	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("ConfirmOrder bad status code:%d", resp.StatusCode)
	}
	body := getBody(resp.Body)
	log.MyOrderLogD("ConfirmOrder:body[%s]", string(body))*/
	return true, nil
}
