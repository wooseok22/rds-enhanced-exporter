package config

import "time"

type ConfInfo struct {
	Global struct {
		AwsArn      string
		Region      string
		Port        string
		SInterval   time.Duration
		TInterval   time.Duration
		Log         string
		SentryDSN   string
		ScrapMethod string
		TSDBUrl     string
	}

	Labels struct {
		Target string
		Kv     []string
	}
}
