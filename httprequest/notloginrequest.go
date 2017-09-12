package httprequest

import (
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

func LeftTicket(client *http.Client, date, fromStation, toStation, code string) (LeftTicketsMsgData, error) {
	resp, err := client.Get(getLeftTicketUrl(date, fromStation, toStation, code))
	if err != nil {
		return LeftTicketsMsgData{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return LeftTicketsMsgData{}, fmt.Errorf("LeftTicket bad status code:%d", resp.StatusCode)
	}

	body := getBody(resp.Body)
	if len(body) == 0 {
		return LeftTicketsMsgData{}, fmt.Errorf("LeftTicket data is empty")
	}
	ltm := LeftTicketsMsg{}
	errJson := json.Unmarshal(body, &ltm)
	fmt.Printf("body:%s\n", string(body))
	return ltm.Data, fmt.Errorf("json parse err:%v", errJson)
}
