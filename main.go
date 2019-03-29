package main

import "sync"

func main() {
	loadConfig()
	cookies := getCookiesFromFile()
	acceptMap := getSubmissions(cookies)

	wg := sync.WaitGroup{}
	wg.Add(len(acceptMap))

	for _, problem := range acceptMap {
		go getCode(cookies, problem, &wg)
	}
	wg.Wait()
}
