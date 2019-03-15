package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"jirabot/jiraapi"
	"log"
	"reflect"
	"regexp"
	"time"
)

// конфиг микросервиса
var config struct {
	BotToken      string
	Jira          jiraapi.Config
	ChatWhiteList []int64
}

// инициализация конфига
func Init(configPath string) {
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		logrus.Fatalln("Не удалось загрузить конфиг", err)
	}
	jiraapi.Init(config.Jira)
}

func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func validateChat(message *tb.Message) bool {
	var chatid int64
	chatid = message.Chat.ID
	valid, _ := in_array(chatid, config.ChatWhiteList)

	logrus.Infof("%v", valid, message.Chat.ID, config.ChatWhiteList)
	return valid
}

func main() {

	Init("config.toml")

	var issueName = regexp.MustCompile(`[A-Z]{2,5}-\d*`)

	b, err := tb.NewBot(tb.Settings{
		Token:  config.BotToken,
		URL:    "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tb.OnText, func(m *tb.Message) {
		if !validateChat(m) {
			return
		}

		var issue string
		issue = issueName.FindString(m.Text)

		if issue != "" {
			text := jiraapi.GetIssueLink(issue)
			_, err := b.Send(m.Chat, text, tb.ModeHTML)

			if err != nil {
				fmt.Printf("%v", err)
			}
		}
	})

	b.Start()
}