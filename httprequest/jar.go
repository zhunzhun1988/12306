package httprequest

import (
	"12306/log"
	"net/http"
	"net/url"
	"sync"
	//"github.com/golang/glog"
)

type jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func newJar() *jar {
	jar := new(jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()

	for _, c := range cookies {

		log.MyLogDebug("cookies[%s/%s]:%s:%s\n", u.Host, c.Path, c.Name, c.Value)
	}
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}
