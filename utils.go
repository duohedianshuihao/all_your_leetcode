package main

import (
	"github.com/spf13/viper"
)

// UserInfo to show the leetcode user information
type RespData struct {
	Username string `json:"user_name"`
}

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
		panic(e)
	}
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	check(err)
}
