package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func verifyCookies(cookies []*http.Cookie) bool {
	userInfo := viper.GetStringMapString("UserInfo")
	URL := viper.GetStringMapString("URL")
	verifyURL := URL["baseURL"] + URL["verifyURL"]

	client := &http.Client{}
	req, err := http.NewRequest("GET", verifyURL, nil)
	check(err)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	check(err)

	respData := RespData{}
	json.NewDecoder(resp.Body).Decode(&respData)

	return respData.Username != userInfo["username"]
}

func getSubmissions(cookies []*http.Cookie) map[string]Problem {
	URL := viper.GetStringMapString("URL")
	language := viper.GetStringMapString("Language")

	submissionURL := URL["baseURL"] + URL["submissionURL"]

	client := &http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", submissionURL, nil)
	check(err)

	query := req.URL.Query()
	limit := 20
	offset := 0

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	submission := &SubmissionResp{
		HasNext: true,
		Lastkey: "",
	}
	acceptMap := make(map[string]Problem)

	for submission.HasNext == true {
		buildURL(req, query, &offset, &limit, submission.Lastkey)
		resp, err := client.Do(req)
		check(err)

		// sleep 0.5 second to ensure the resp is done
		time.Sleep(500 * time.Millisecond)
		json.NewDecoder(resp.Body).Decode(submission)
		updateProblem(acceptMap, submission.ProblemsDump, language["lang"])
	}
	return acceptMap
}

func getCode(cookies []*http.Cookie, problem Problem) string {
	URL := viper.GetStringMapString("URL")
	language := viper.GetStringMapString("Language")

	targetURL := URL["baseURL"] + problem.URL

	client := &http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", targetURL, nil)
	check(err)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	check(err)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	role := regexp.MustCompile(`submissionCode: \'(.*)\'`)
	code, err := strconv.Unquote("\"" + role.FindStringSubmatch(bodyString)[1] + "\"")

	errPath := os.Chdir(language["download"])
	check(errPath)
	file, err := os.Create(problem.Title + "." + language["postfix"])
	check(err)

	defer file.Close()

	file.WriteString(code)

	return code
}

func updateProblem(acceptMap map[string]Problem, problemDump []Problem, lang string) {

	for _, problem := range problemDump {

		if problem.StatusDisplay != "Accepted" || problem.Lang != lang {
			continue
		}

		if val, ok := acceptMap[problem.Title]; ok {
			if compareRuntime(problem.Runtime, val.Runtime) {
				acceptMap[problem.Title] = problem
			}
			continue
		}
		acceptMap[problem.Title] = problem
	}
}

func buildURL(req *http.Request, query url.Values, offset *int, limit *int, lastkey string) {
	query.Set("lastkey", lastkey)
	query.Set("limit", strconv.Itoa(*limit))
	query.Set("offset", strconv.Itoa(*offset))

	*offset += *limit
	req.URL.RawQuery = query.Encode()
}

func compareRuntime(time1 string, time2 string) bool {
	re := regexp.MustCompile(`[0-9]*`)
	t1, err1 := strconv.Atoi(re.FindString(time1))
	check(err1)
	t2, err2 := strconv.Atoi(re.FindString(time2))
	check(err2)
	return t1 < t2
}
