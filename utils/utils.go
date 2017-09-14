package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var numRand *rand.Rand = nil

func getNumRand() *rand.Rand {
	if numRand == nil {
		numRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return numRand
}

func GetNowMicoSecondStr() string {
	n := time.Now()
	return fmt.Sprintf("%d", n.Unix()*int64(1000)+int64(n.Nanosecond()/1000000))
}

func GetRandFloat(len int) string {
	r := getNumRand()
	tmp := "0."
	for i := 0; i < len; i++ {
		tmp += strconv.Itoa(r.Int() % 10)
	}
	return tmp
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

func WriteFile(filename string, reader io.ReadCloser) error {
	var f *os.File
	var err error
	if CheckFileIsExist(filename) { //如果文件存在
		//f, err = os.OpenFile(filename, os.O_APPEND, 0777) //打开文件
		errDel := os.Remove(filename)
		if errDel != nil {
			return errDel
		}
	}
	f, err = os.Create(filename) //创建文件

	if err != nil {
		return err
	}
	_, err = io.Copy(f, reader)
	return err
}

func GetCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
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

func GetOrderTimeFomat(t time.Time, isDay bool) string {
	if isDay {
		return fmt.Sprintf("%s %s %02d %d 00:00:00 GMT+0800 (中国标准时间)", t.Weekday().String()[0:3], t.Month().String()[0:3],
			t.Day(), t.Year())
	}

	return fmt.Sprintf("%s %s %02d %d %02d:%02d:%02d GMT+0800 (中国标准时间)", t.Weekday().String()[0:3], t.Month().String()[0:3],
		t.Day(), t.Year(), t.Hour(), t.Minute(), t.Second())
}

func UrlEncode(dataStr string) string {
	data := []byte(dataStr)
	t := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}
	ret := make([]byte, 0, len(data)*2)
	for _, b := range data {
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') {
			ret = append(ret, b)
		} else {
			ret = append(ret, '%')
			ret = append(ret, t[int(b)/16])
			ret = append(ret, t[int(b)%16])
		}
	}
	return string(ret)
}
