package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

func leetcodeLogin() []*http.Cookie {
	userInfo := viper.GetStringMapString("userinfo")

	URL := viper.GetStringMapString("leetcode")
	loginURL := URL["baseurl"] + URL["loginurl"]

	client := &http.Client{}

	resp, err := client.Get(URL["baseurl"])
	check(err)

	cookies := resp.Cookies()

	data := url.Values{
		"csrfmiddlewaretoken": {resp.Cookies()[1].Value},
		"login":               {userInfo["login"]},
		"password":            {userInfo["password"]},
		"next":                {"/problems"},
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	req.Header.Set("referer", loginURL)
	req.Header.Set("origin", URL["baseurl"])
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	loginResp, err := client.Do(req)
	check(err)

	if loginResp.StatusCode != http.StatusOK {
		panic("login failed")
	}

	LCookies := loginResp.Cookies()

	// writeFile("./", "cookies", resp.Header["Set-Cookie"])

	return LCookies

}

func getCookiesFromFile() []*http.Cookie {

	dat, err := ioutil.ReadFile("cookies")
	check(err)
	data := strings.Split(string(dat), ";")
	cookies := make([]*http.Cookie, len(data))
	for i, val := range data {
		val := strings.TrimSpace(val)
		cookies[i] = &http.Cookie{
			Name:  strings.Split(val, "=")[0],
			Value: strings.Split(val, "=")[1],
		}

	}
	return cookies
}
