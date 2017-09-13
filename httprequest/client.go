package httprequest

import (
	"12306/log"
	"12306/utils"
	"12306/verifycode"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"time"
)

const (
	request_timeout = 30 * time.Second
)

type Interface interface {
	Login() error
	IsLogined() bool
	GetPassengers() ([]Passenger, error)
	GetStations() ([]StationItem, error)
	GetLeftTickets(date, fromStation, toStation string) (TicketsInfoList, error)
}
type Client struct {
	client             *http.Client
	isLogined          bool
	verifies           verifycode.VerifierList
	stationCache       []StationItem
	username, password string
}

func NewClient(username, password string) Interface {
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
		isLogined:    false,
		stationCache: nil,
		username:     username,
		password:     password,
		verifies:     verifycode.VerifierList{verifycode.NewDebugVerify()},
	}
}

func getNextFileName() string {
	dir := path.Join(utils.GetCurrentPath(), "image")
	files, err := utils.ListDir(dir, "image_")
	if err != nil {
		if utils.CheckFileIsExist(dir) == false {
			os.Mkdir(dir, 0777)
		}
		return path.Join(dir, "image_0001.jpg")
	} else {
		return path.Join(dir, fmt.Sprintf("image_%04d.jpg", len(files)+1))
	}
}
func (c *Client) Login() error {
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
	errWebLogin := WebLogin(c.client, c.username, c.password)
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
	return UserLoginCheck(c.client)
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

func (c *Client) GetStations() ([]StationItem, error) {
	if c.stationCache != nil && len(c.stationCache) > 0 {
		return c.stationCache, nil
	}
	log.MyLogInfo("获取车站信息...")
	ret, err := GetStations(c.client)
	if err == nil {
		c.stationCache = ret
	}
	log.MyLogInfo("车站信息数:%d", len(ret))
	return ret, err
}

func changeStationNameToCode(stations []StationItem, name string) string {
	for _, s := range stations {
		if s.Name == name {
			return s.ID
		}
	}
	return ""
}

func (c *Client) GetLeftTickets(date, fromStation, toStation string) (TicketsInfoList, error) {
	_, err := c.GetStations()
	if err != nil {
		return TicketsInfoList{}, fmt.Errorf("GetLeftTickets query station info err:%v", err)
	}
	from, to := changeStationNameToCode(c.stationCache, fromStation), changeStationNameToCode(c.stationCache, toStation)
	if date == "" || from == "" || to == "" {
		return TicketsInfoList{}, fmt.Errorf("GetLeftTickets parms err")
	}
	return LeftTicket(c.client, date, from, to, "ADULT")
}
