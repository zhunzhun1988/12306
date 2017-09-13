package httprequest

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
}

type LeftTicketsMsgData struct {
	Result []string          `json:"result"`
	Flag   string            `json:"flag"`
	Map    map[string]string `json:"map"`
}

type LeftTicketsMsg struct {
	ValidateMessagesShowId string             `json:"validateMessagesShowId"`
	Status                 string             `json:"status"`
	Httpstatus             int                `json:"httpstatus"`
	Data                   LeftTicketsMsgData `json:"Data"`
	Messages               string             `json:"messages"`
	ValidateMessages       string             `json:"validateMessages"`
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
