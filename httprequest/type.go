package httprequest

import (
	"fmt"
	"strings"
	"time"
)

const (
	Verify_OK_CODE = "4"
	Login_OK_Code  = 0
)

type VerifyMessage struct {
	Result_message string `json:"result_message"`
	Result_code    string `json:"result_code"`
}

type LoginMessage struct {
	Result_message string `json:"result_message"`
	Result_code    int    `json:"result_code"`
	Uamtk          string `json:"uamtk"`
}

type TokenMessage struct {
	Result_message string `json:"result_message"`
	Result_code    int    `json:"result_code"`
	AppTK          string `json:"apptk"`
	NewAppTK       string `json:"newapptk"`
}

type StationItem struct {
	PYJianXie string
	Name      string
	ID        string
	PingYin   string
	Code      string
	Index     string
}

type TicketsInfo struct {
	SecretStr   string
	FromStation string
	ToStation   string
	Message     string
	StartTime   time.Time
	EndTime     time.Time
	TrianName   string
	HaveTickets bool
	TDZ         string
	YDZ         string
	EDZ         string
	GJRW        string
	RW          string
	DW          string
	YW          string
	RZ          string
	YZ          string
	WZ          string
	QT          string
}

func (t TicketsInfo) ToString() string {
	strF := func(str string) string {
		if str != "" {
			return str
		}
		return "--"
	}
	return fmt.Sprintf("车次：[%s]\t出发：[%s]\t到达:[%s]\t\t日期：[%s]\t特等坐:[%s]\t一等座:[%s]\t二等座:[%s]\t高级软卧:[%s]\t软卧:[%s]\t动卧:[%s]\t硬卧:[%s]\t软座:[%s]\t硬座:[%s]\t无座:[%s]",
		t.TrianName, t.FromStation, t.ToStation, t.StartTime.Format("2006-01-02"), strF(t.TDZ), strF(t.YDZ), strF(t.EDZ), strF(t.GJRW), strF(t.RW),
		strF(t.DW), strF(t.YW), strF(t.RZ), strF(t.YZ), strF(t.WZ))
}

type TicketsInfoList []TicketsInfo

func (tl TicketsInfoList) ToStrings() []string {
	ret := make([]string, 0, len(tl)+1)
	strF := func(str string) string {
		if str != "" {
			return str
		}
		return "--"
	}
	ret = append(ret, fmt.Sprintf("车次\t出发\t到达\t特等坐\t一等座\t二等座\t高级软卧\t软卧\t动卧\t硬卧\t软座\t硬座\t无座\t其他\t备注\t有票"))
	for _, t := range tl {
		ret = append(ret, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%t",
			t.TrianName, t.FromStation, t.ToStation, strF(t.TDZ), strF(t.YDZ), strF(t.EDZ),
			strF(t.GJRW), strF(t.RW), strF(t.DW), strF(t.YW), strF(t.RZ), strF(t.YZ), strF(t.WZ), strF(t.QT), t.Message, t.HaveTickets))
	}
	return ret
}

func stringToTicketsInfo(str string, stationMap map[string]string) TicketsInfo {
	ret := TicketsInfo{}
	strs := strings.Split(str, "|")
	if len(strs) != 36 {
		return ret
	}
	ret.SecretStr = strings.Replace(strings.Replace(strings.Replace(strings.Replace(strs[0], "%2F", "/", -1), "%2B", "+", -1), "%3D", "=", -1), "%0A", " ", -1)
	ret.Message = strs[1]
	ret.TrianName = strs[3]

	ret.FromStation = stationMap[strs[6]]
	ret.ToStation = stationMap[strs[7]]
	ret.StartTime, _ = time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-%s %s:00", string(strs[13][0:4]), string(strs[13][4:6]), string(strs[13][6:]), strs[8]))
	ret.EndTime, _ = time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-%s %s:00", string(strs[13][0:4]), string(strs[13][4:6]), string(strs[13][6:]), strs[9]))
	if strs[11] == "Y" {
		ret.HaveTickets = true
	}
	ret.DW = strs[33]  // ok
	ret.TDZ = strs[32] // ok
	ret.YDZ = strs[31] // ok
	ret.EDZ = strs[30] // ok
	ret.YZ = strs[29]  // ok
	ret.YW = strs[28]  // ok
	ret.WZ = strs[26]  // ok
	ret.QT = strs[25]
	ret.RZ = strs[24]
	ret.RW = strs[23]   //ok
	ret.GJRW = strs[21] //ok
	return ret
}

func LeftTicketsMsgDataToTicketsInfoList(lt *LeftTicketsMsgData) TicketsInfoList {
	ret := make(TicketsInfoList, 0, len(lt.Result))
	for _, str := range lt.Result {
		ret = append(ret, stringToTicketsInfo(str, lt.Map))
	}
	return ret
}

type LeftTicketsMsgData struct {
	Result []string          `json:"result"`
	Flag   string            `json:"flag"`
	Map    map[string]string `json:"map"`
}

type LeftTicketsMsg struct {
	ValidateMessagesShowId string             `json:"validateMessagesShowId"`
	Status                 bool               `json:"status"`
	Httpstatus             int                `json:"httpstatus"`
	Data                   LeftTicketsMsgData `json:"Data"`
	//Messages               string             `json:"messages"`
	//ValidateMessages       string             `json:"validateMessages"`
}

type LoginCheckMsgData struct {
	Flag bool `json:"flag"`
}
type LoginCheckMsg struct {
	ValidateMessagesShowId string            `json:"validateMessagesShowId"`
	Status                 bool              `json:"status"`
	Httpstatus             int               `json:"httpstatus"`
	Data                   LoginCheckMsgData `json:"Data"`
	//Messages               []string          `json:"messages"`
	//ValidateMessages string `json:"validateMessages"`
}

type LeftTicketLoginDeviceMsg struct {
	Exp string `json:"exp"`
	Dfp string `json:"dfp"`
}

type OrderTicketMsg struct {
	ValidateMessagesShowId string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	Httpstatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
}
