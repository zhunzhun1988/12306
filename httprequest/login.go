package httprequest

import (
	"12306/log"
	"12306/utils"
	"12306/verifycode"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
)

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

func GetLoginVerifyImg(client *http.Client, imageSavePath string) error {
	resp, err := client.Get(getLoginVerifyImgUrl())
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GetLoginVerifyImg bad status code:%d", resp.StatusCode)
	}
	return utils.WriteFile(imageSavePath, resp.Body)
}

func LoginInit(client *http.Client) error {
	resp, err := client.Get(login_init_addr)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LoginInit bad status code:%d", resp.StatusCode)
	}

	/*	resp, err = client.Get("https://kyfw.12306.cn/otn/dynamicJs/lpkfrls")
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("LoginInit bad status code:%d", resp.StatusCode)
		}*/
	return nil
}

func LoginInit12306(client *http.Client) error {
	resp, err := client.Get(login_init12306_addr)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LoginInit12306 bad status code:%d", resp.StatusCode)
	}
	return nil
}

func CheckVerifiyLoginCode(client *http.Client, poss verifycode.VerifyPosList) error {
	resp, err := client.PostForm(login_verify_addr, getLoginVerifyUrlValues(poss))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("VerifiyLoginCode bad status code:%d", resp.StatusCode)
	}

	var vm VerifyMessage = VerifyMessage{}
	errJson := json.Unmarshal(getBody(resp.Body), &vm)
	if errJson != nil {
		return fmt.Errorf("VerifiyLoginCode Unmarshal error:%d", errJson)
	}
	if vm.Result_code != Verify_OK_CODE {
		return fmt.Errorf("VerifiyLoginCode:%v", vm)
	}
	return nil
}

func WebLogin(client *http.Client, username, password string) error {
	resp, err := client.PostForm(weblogin_addr, getLoginUrlValues(username, password))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WebLogin bad status code:%d", resp.StatusCode)
	}
	var lm LoginMessage = LoginMessage{}
	errJson := json.Unmarshal(getBody(resp.Body), &lm)
	if errJson != nil {
		return fmt.Errorf("WebLogin Unmarshal error:%v", errJson)
	}
	log.MyLogDebug("loginMsg:%v", lm)
	if lm.Result_code != Login_OK_Code {
		return fmt.Errorf("WebLogin:%v", lm)
	}
	return nil
}

func AuthUamtk(client *http.Client) error {
	resp, err := client.PostForm(get_token_addr, url.Values{"appid": []string{"otn"}})
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AuthUamtk bad status code:%d", resp.StatusCode)
	}
	buf := getBody(resp.Body)
	log.MyLogDebug("uamtk body:%s", buf)

	var tm TokenMessage = TokenMessage{}
	errJson := json.Unmarshal([]byte(buf), &tm)
	if errJson != nil {
		return fmt.Errorf("AuthUamtk Unmarshal error:%s", string(buf))
	}
	if tm.NewAppTK == "" {
		return fmt.Errorf("AuthUamtk get token error:%s", string(buf))
	}
	log.MyLogDebug("get new token:%s", tm.NewAppTK)
	resp, err = client.PostForm(set_token_addr, url.Values{"tk": []string{tm.NewAppTK}})
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("uamauthclient bad status code:%d", resp.StatusCode)
	}
	log.MyLogDebug("uamauthclient body:%s", getBody(resp.Body))
	return nil
}

func getBody(r io.ReadCloser) []byte {
	buf, errRead := ioutil.ReadAll(r)
	if errRead != nil {
		return []byte{}
	}
	return buf

}
func UserLogin(client *http.Client) error {
	log.MyLogDebug("user login step1")
	resp, err := client.PostForm(userlogin_addr1, url.Values{"_json_att": []string{}})
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("UserLogin1 bad status code:%d", resp.StatusCode)
	}

	log.MyLogDebug("user login step2")
	resp, err = client.Get(userlogin_addr2)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("UserLogin2 bad status code:%d", resp.StatusCode)
	}

	return nil
}
