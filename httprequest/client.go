package httprequest

import (
	"12306/log"
	"12306/verifycode"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

const (
	request_timeout = 30 * time.Second
)

type Client struct {
	client    *http.Client
	isLogined bool
	verifies  verifycode.VerifierList
}

func NewClient() *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, request_timeout) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(request_timeout)) //设置发送接受数据超时
			return conn, nil
		},
	}
	jar, _ := cookiejar.New(nil)
	return &Client{
		client: &http.Client{
			Jar:       jar, //newJar(),
			Transport: tr,
		},
		isLogined: false,
		verifies:  verifycode.VerifierList{verifycode.NewDebugVerify()},
	}
}

func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasPrefix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

func getNextFileName() string {
	files, err := ListDir("/home/adam/go/src/12306/image", "image_")
	if err != nil {
		return "image/image_0001.jpg"
	} else {
		return fmt.Sprintf("image/image_%04d.jpg", len(files)+1)
	}
}
func (c *Client) Login(username, password string) error {
	log.MyLoginLogI("开始登录...")
	errInit := LoginInit(c.client)
	if errInit != nil {
		log.MyLoginLogE("登录失败：%v\n", errInit)
		return errInit
	}

	log.MyLogDebug("开始拉取登录验证码")
	saveFile := getNextFileName()
	errVerify := GetLoginVerifyImg(c.client, saveFile)
	if errVerify != nil {
		log.MyLoginLogE("登录失败：%v\n", errVerify)
		return errVerify
	}
	log.MyLogDebug("开始验证验证码")
	errCheck := CheckVerifiyLoginCode(c.client, c.verifies.GetAnswer(saveFile))
	if errCheck != nil {
		log.MyLoginLogE("登录失败：%v", errCheck)
		return errCheck
	}
	log.MyLogDebug("开始用户登录")
	errWebLogin := WebLogin(c.client, username, password)
	if errWebLogin != nil {
		log.MyLoginLogE("登录失败：%v", errCheck)
		return errWebLogin
	}
	log.MyLogDebug("开始正式用户登录")
	errUserLogin := UserLogin(c.client)
	if errUserLogin != nil {
		log.MyLoginLogE("登录失败：%v", errUserLogin)
		return errUserLogin
	}

	log.MyLogDebug("开始获取token")
	authErr := AuthUamtk(c.client)
	if authErr != nil {
		log.MyLoginLogE("登录失败：%v", authErr)
		return authErr
	}

	log.MyLogDebug("模拟12306跳转")
	errInit12306 := LoginInit12306(c.client)
	if errInit12306 != nil {
		log.MyLoginLogE("登录失败：%v\n", errInit12306)
		return errInit12306
	}
	log.MyLoginLogI("登录成功")
	c.isLogined = true
	return nil
}

func (c *Client) IsLogined() bool {
	if c.isLogined == false {
		return false
	}
	return true
}
func (c *Client) GetPassengers() ([]Passenger, error) {
	if c.isLogined == false {
		log.MyLog(log.ERROR, log.PASSENGER, "获取用户信息失败:未登录")
		return []Passenger{}, fmt.Errorf("请先登录")
	}
	log.MyLog(log.INFO, log.PASSENGER, "获取用户信息...")
	ps, err := GetPassengers(c.client)
	if err != nil {
		log.MyLog(log.ERROR, log.PASSENGER, "获取用户信息失败:%s", err)
	}
	log.MyLog(log.DEBUG, log.PASSENGER, "用户信息[%v]", ps)
	return ps, err
}
