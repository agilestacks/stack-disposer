package config

import "time"

var (
	Port    string
	GitUrl  string
	GitDir  string
	Verbose bool
	Timeout time.Duration
)
