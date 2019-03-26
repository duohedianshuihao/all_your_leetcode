package main

func main() {
	loadConfig()
	cookies := getCookiesFromFile()
	acceptMap := getSubmissions(cookies)

	for _, problem := range acceptMap {
		go getCode(cookies, problem)
	}
}
