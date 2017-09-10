package httprequest

import (
	"12306/log"
	"12306/utils"
	"12306/verifycode"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
)

const (
	login_verify_codeimage_addr = "https://kyfw.12306.cn/passport/captcha/captcha-image?login_site=E&module=login&rand=sjrand&"
	login_verify_addr           = "https://kyfw.12306.cn/passport/captcha/captcha-check"
	weblogin_addr               = "https://kyfw.12306.cn/passport/web/login"
	userlogin_addr1             = "https://kyfw.12306.cn/otn/login/userLogin"
	userlogin_addr2             = "https://kyfw.12306.cn/otn/passport?redirect=/otn/login/userLogin"
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

func CheckVerifiyLoginCode(client *http.Client, poss verifycode.VerifyPosList) error {
	resp, err := client.PostForm(login_verify_addr, getLoginVerifyUrlValues(poss))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("VerifiyLoginCode bad status code:%d", resp.StatusCode)
	}
	buf, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return fmt.Errorf("VerifiyLoginCode read error:%d", errRead)
	}

	if string(buf) != "" {
		var vm VerifyMessage = VerifyMessage{}
		errJson := json.Unmarshal(buf, &vm)
		if errJson != nil {
			return fmt.Errorf("VerifiyLoginCode Unmarshal error:%d", errJson)
		}
		if vm.Result_code != Verify_OK_CODE {
			return fmt.Errorf("VerifiyLoginCode:%v", vm)
		}
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
	buf, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return fmt.Errorf("WebLogin read error:%d", errRead)
	}

	if string(buf) != "" {
		var lm LoginMessage = LoginMessage{}
		errJson := json.Unmarshal(buf, &lm)
		if errJson != nil {
			return fmt.Errorf("WebLogin Unmarshal error:%s", string(buf))
		}
		log.MyLogDebug("loginMsg:%v", lm)
		if lm.Result_code != Login_OK_Code {
			return fmt.Errorf("WebLogin:%v", lm)
		}
	}
	return nil
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

	log.MyLogDebug("user login step3")
	resp, err = client.PostForm("https://kyfw.12306.cn/passport/web/auth/uamtk", url.Values{"appid": []string{"otn"}})
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("UserLogin3 bad status code:%d", resp.StatusCode)
	}

	log.MyLogDebug("user login step4")
	resp, err = client.PostForm("https://kyfw.12306.cn/otn/uamauthclient", url.Values{"tk": []string{"zvaAMSVnBTpijXCC6IGNrrUjap17wXW_85MuzMTW6qfOjRQsfs6260"}})
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("UserLogin4 bad status code:%d", resp.StatusCode)
	}
	/*buf, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return fmt.Errorf("UserLogin read error:%d", errRead)
	}

	if string(buf) != "" {
		var lm LoginMessage = LoginMessage{}
		errJson := json.Unmarshal(buf, &lm)
		if errJson != nil {
			return fmt.Errorf("UserLogin Unmarshal error:%s", string(buf))
		}
		log.MyLogDebug("UserLogin:%v", lm)
		if lm.Result_code != Login_OK_Code {
			return fmt.Errorf("StartLogin:%v", lm)
		}
	}*/
	return nil
}
