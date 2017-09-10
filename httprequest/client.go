package httprequest

import (
	"12306/log"
	"12306/verifycode"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
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
	return &Client{
		client: &http.Client{
			Jar:       newJar(),
			Transport: tr,
		},
		isLogined: false,
		verifies:  verifycode.VerifierList{verifycode.NewDebugVerify()},
	}
}

func (c *Client) Login(username, password string) error {
	log.MyLoginLogI("开始登录...")
	errVerify := GetLoginVerifyImg(c.client, "loginverifycode.jpg")
	if errVerify != nil {
		log.MyLoginLogE("登录失败：%s\n", errVerify)
		return errVerify
	}

	errCheck := CheckVerifiyLoginCode(c.client, c.verifies.GetAnswer("loginverifycode.jpg"))
	if errCheck != nil {
		log.MyLoginLogE("登录失败：%s", errCheck)
		return errCheck
	}

	errWebLogin := WebLogin(c.client, username, password)
	if errWebLogin != nil {
		log.MyLoginLogE("登录失败：%s", errCheck)
		return errWebLogin
	}

	errUserLogin := UserLogin(c.client)
	if errUserLogin != nil {
		log.MyLoginLogE("登录失败：%s", errUserLogin)
		return errUserLogin
	} else {
		log.MyLoginLogI("登录成功")
	}
	c.isLogined = true
	return nil
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
