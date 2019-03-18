package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	bp "jirabot/bot"
	conf "jirabot/config"
	"jirabot/handler"
	"jirabot/jiraapi"
	"ms-acr/acr/consts"
	"net/http"
	"os"
)

const version = 1

// конфиг микросервиса
var config conf.Config

// инициализация конфига
func Init(configPath string) {
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		logrus.Fatalln("Не удалось загрузить конфиг", err)
	}
	jiraapi.Init(config.Jira)
}

func main() {
	a := cli.NewApp()

	a.Name = "jirabot"
	a.Usage = "Telegram bot for Jira Atlassian Stack"
	a.Version = consts.APP_VERSION
	a.Author = "Nikolay Kindyakov"
	a.Email = "kindyakov.nikolay@kolesa.kz"
	a.Action = actionRun
	a.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, b",
			Usage: "If provided, the service will be launched in debug mode",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: "/etc/telejira/config.cfg",
			Usage: "Path to the configuration file",
		},
	}
	a.Run(os.Args)
}

func actionRun(c *cli.Context) {
	var updates tgbotapi.UpdatesChannel

	file := c.String("config")
	isDebug := c.Bool("debug")

	Init(file)

	botApi, err := tgbotapi.NewBotAPI(config.Bot.Token)

	if err != nil {
		fmt.Println(err)
	}

	if isDebug {
		botApi.Debug = true
	} else {
		botApi.Debug = false
	}

	if config.Bot.UseWebHook {
		updates = fetchWebhookUpdates(botApi)
	} else {
		updates = fetchPollingUpdates(botApi)
	}

	runServer()

	bot := bp.BotApi{botApi, updates, tgbotapi.Update{}}

	fmt.Println("Listening to bot api updates")

	bot.ListenForUpdates(&config)
}

func fetchPollingUpdates(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	res, _ := bot.RemoveWebhook()
	fmt.Printf("%f", res)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = config.Bot.PollingTimeout

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		fmt.Errorf("Problem in setting Polling", err.Error())
	}

	return updates
}

func fetchWebhookUpdates(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	_, err := bot.SetWebhook(tgbotapi.NewWebhook(config.Bot.WebHook + "/bot" + bot.Token))

	if err != nil {
		fmt.Errorf("Problem in setting Webhook", err.Error())
	}

	updates := bot.ListenForWebhook("/bot" + bot.Token)

	return updates
}

func runServer() {
	http.HandleFunc("/", handler.MainHandler)
	http.HandleFunc("/health", handler.HealthHandler)
	go http.ListenAndServe(config.Bot.HttpPort, nil)

	fmt.Println("Starting server at port " + config.Bot.HttpPort)
}
