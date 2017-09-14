package httprequest

import (
	"12306/log"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func parseStation(str string) []StationItem {
	strs := strings.Split(str, "@")
	ret := make([]StationItem, 0, len(strs))
	for _, itemStr := range strs {
		if itemStr == "" {
			continue
		}
		items := strings.Split(itemStr, "|")
		if len(items) < 6 {
			continue
		}
		ret = append(ret, StationItem{
			PYJianXie: items[0],
			Name:      items[1],
			ID:        items[2],
			PingYin:   items[3],
			Code:      items[4],
			Index:     items[5],
		})
	}
	return ret
}
func GetStations(client *http.Client) ([]StationItem, error) {
	ret := []StationItem{}
	resp, err := client.Get(get_station_addr)
	if err != nil {
		return ret, err
	}
	if resp.StatusCode != http.StatusOK {
		return ret, fmt.Errorf("GetStations bad status code:%d", resp.StatusCode)
	}

	body := string(getBody(resp.Body))
	if body == "" || strings.HasPrefix(body, "var station_names =") == false {
		return ret, fmt.Errorf("GetStations unknow data:%s", body)
	}

	start := strings.Index(body, "'")
	end := strings.LastIndex(body, "'")
	if start <= 0 || end <= start {
		return ret, fmt.Errorf("GetStations unknow data:%s", body)
	}
	tmp := body[start+1 : end]
	return parseStation(tmp), nil
}

var exp string = ""
var dfp string = ""

func getExpAndDfp(client *http.Client) (expret, dfpret string) {
	if exp != "" && dfp != "" {
		return exp, dfp
	}
	resp, err := client.Get(getLeftTicketLoginDeviceUrl())
	if err != nil {
		return "", ""
	}
	body := string(getBody(resp.Body))
	if body == "" || strings.HasPrefix(body, "callbackFunction") == false {
		return "", ""
	}
	startIndex := strings.Index(body, "'")
	endIndex := strings.LastIndex(body, "'")
	if startIndex <= 0 || endIndex <= startIndex {
		return "", ""
	}
	jsonStr := body[startIndex+1 : endIndex]
	lldm := LeftTicketLoginDeviceMsg{}
	errJson := json.Unmarshal([]byte(jsonStr), &lldm)
	if errJson != nil {
		return "", ""
	}
	exp = lldm.Exp
	dfp = lldm.Dfp
	return exp, dfp
}

func LeftTicket(client *http.Client, date, fromStation, toStation, code string) (TicketsInfoList, error) {
	curExp, curDfp := getExpAndDfp(client)
	req, _ := http.NewRequest("Get", getLeftTicketUrl(date, fromStation, toStation, code), nil)
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	req.AddCookie(&http.Cookie{Name: "RAIL_EXPIRATION", Value: curExp})
	req.AddCookie(&http.Cookie{Name: "RAIL_DEVICEID", Value: curDfp})

	resp, err := client.Do(req)
	if err != nil {
		return TicketsInfoList{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return TicketsInfoList{}, fmt.Errorf("LeftTicket bad status code:%d", resp.StatusCode)
	}

	body := getBody(resp.Body)
	log.MyLogDebug("LeftTicket body:[%s]", string(body))
	if len(body) == 0 {
		return TicketsInfoList{}, fmt.Errorf("LeftTicket data is empty")
	}
	ltm := LeftTicketsMsg{}
	errJson := json.Unmarshal(body, &ltm)
	if errJson != nil {
		return TicketsInfoList{}, fmt.Errorf("json parse err:%v, [%s]", errJson, string(body))
	}
	return LeftTicketsMsgDataToTicketsInfoList(&ltm.Data), nil
}
