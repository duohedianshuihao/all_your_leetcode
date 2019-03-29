package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type SubmissionResp struct {
	Lastkey      string    `json:"last_key"`
	HasNext      bool      `json:"has_next"`
	ProblemsDump []Problem `json:"submissions_dump"`
}

type Problem struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	Lang          string `json:"lang"`
	StatusDisplay string `json:"status_display"`
	Runtime       string `json:"runtime"`
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	// viper.SetKeyCaseSensitivity(true)
	err := viper.ReadInConfig()
	check(err)
}

func writeFile(path string, filename string, content string) {
	os.Chdir(path)

	fd, err := os.Create(filename)
	check(err)
	defer fd.Close()

	fd.WriteString(content)

}
