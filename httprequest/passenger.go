package httprequest

import (
	"12306/log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	get_passenger_addr = "https://kyfw.12306.cn/otn/passengers/init"
)

type Passenger struct {
	PassengerTypeName   string `json:"passenger_type_name"`
	IsMyself            string `json:"isUserSelf"`
	PassengerIDTypeCode string `json:"passenger_id_type_code"`
	PassengerName       string `json:"passenger_name"`
	TotalTime           string `json:"total_times"`
	PassengerIDTypeName string `json:"passenger_id_type_name"`
	PassengerType       string `json:"passenger_type"`
	PassengerIDNo       string `json:"passenger_id_no"`
	MobileNo            string `json:"mobile_no"`
}

func getPassengerUrlValues() url.Values {
	ret := make(url.Values)
	ret["_json_att"] = []string{}
	return ret
}

func getPassengerStr(buf string) string {
	strs := strings.Split(buf, "\n")
	passengerKey := "passengers="
	for _, str := range strs {
		if index := strings.Index(str, passengerKey); index > 0 {
			tmp := str[index+len(passengerKey):]
			return strings.TrimRight(tmp, ";")
		}
	}
	return ""
}

func GetPassengers(client *http.Client) ([]Passenger, error) {
	ret := []Passenger{}
	resp, err := client.PostForm(get_passenger_addr, getPassengerUrlValues())
	if err != nil {
		return ret, err
	}
	if resp.StatusCode != http.StatusOK {
		return ret, fmt.Errorf("GetPassengers bad status code:%d", resp.StatusCode)
	}
	buf, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return ret, fmt.Errorf("GetPassengers read error:%d", errRead)
	}

	if len(buf) != 0 {
		passengerInfo := strings.Replace(getPassengerStr(string(buf)), "'", "\"", -1)
		log.MyLogDebug("passengerInfo:%s", passengerInfo)
		if passengerInfo != "" {
			errJson := json.Unmarshal([]byte(passengerInfo), &ret)
			if errJson != nil {
				return ret, errJson
			}
			return ret, nil
		} else {
			return ret, fmt.Errorf("GetPassengers passengerInfo is empty")
		}
	}
	return ret, fmt.Errorf("GetPassengers read empty")
}
