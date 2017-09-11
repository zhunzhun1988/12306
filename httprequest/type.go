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
