package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

func leetcodeLogin() []*http.Cookie {
	userInfo := viper.GetStringMapString("UserInfo")
	URL := viper.GetStringMapString("Leetcode")
	loginURL := URL["baseURL"] + URL["loginURL"]

	client := &http.Client{}
	resp, err := client.Get(URL["baseURL"])
	check(err)
	fmt.Println(resp.StatusCode)

	cookies := resp.Cookies()

	data := url.Values{
		"csrfmiddlewaretoken": {resp.Cookies()[1].Value},
		"login":               {userInfo["login"]},
		"password":            {userInfo["password"]},
		"next":                {"/problems"},
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	req.Header.Set("referer", loginURL)
	req.Header.Set("origin", URL["baseURL"])
	req.AddCookie(cookies[1])
	loginResp, err := client.Do(req)
	check(err)

	if loginResp.StatusCode != http.StatusOK {
		panic("login failed")
	}

	LCookies := loginResp.Cookies()
	saveCookies(LCookies)

	return LCookies

}

func getCookiesFromFile() []*http.Cookie {

	userInfo := viper.GetStringMapString("UserInfo")

	data := strings.Split(userInfo["cookie"], ";")
	cookies := make([]*http.Cookie, len(data))
	for i, val := range data {
		cookies[i] = &http.Cookie{
			Name:  strings.Split(val, "=")[0],
			Value: strings.Split(val, "=")[1],
		}

	}
	return cookies
}

func saveCookies(cookies []*http.Cookie) {

}
