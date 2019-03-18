package config

import (
	"jirabot/jiraapi"
)

type Config struct {
	Bot  BotConfig
	Jira jiraapi.Config
}

type BotConfig struct {
	Token          string
	WebHook        string
	UseWebHook     bool
	HttpPort       string
	ChatWhiteList  []int64
	PollingTimeout int
}
