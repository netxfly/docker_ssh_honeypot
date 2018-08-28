package models

import "time"

type (
	APIDATA struct {
		Tag      string
		Hostname string
		Content  string
	}

	CheckResult struct {
		Timestamp time.Time
		Ip        string
		Port      string
		Username  string
		Password  string
		Status    int
		Tag       string
	}
)
