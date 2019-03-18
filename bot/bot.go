package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"jirabot/config"
	"jirabot/jiraapi"
	mw "jirabot/middleware"
	"regexp"
)

type BotApi struct {
	BotApi  *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
	Update  tgbotapi.Update
}

func (bot *BotApi) ListenForUpdates(config *config.Config) {
	for update := range bot.Updates {
		bot.Update = update

		if update.Message != nil {
			if !validateChat(update.Message, config) {
				continue
			}

			issueLink := handleIssue(update.Message.Text)

			if issueLink != "" {
				SendTextAnswer(bot, update.Message, issueLink, "Markdown")
			}
		} else {

		}
	}
}

func SendTextAnswer(bot *BotApi, m *tgbotapi.Message, text string, mtype string) {
	if mtype == "" {
		mtype = "Markdown"
	}

	if m.Chat.ID != 0 {
		msg := tgbotapi.NewMessage(m.Chat.ID, text)
		msg.ReplyToMessageID = m.MessageID
		msg.ParseMode = mtype
		bot.BotApi.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(int64(m.From.ID), text)
		bot.BotApi.Send(msg)
	}

}
func handleIssue(m string) string {
	var issue string
	var issueName = regexp.MustCompile(`[A-Za-z]{2,5}-\d*`)
	issue = issueName.FindString(m)

	if issue != "" {
		link := jiraapi.GetIssueLink(issue)

		regexp.MustCompile(`h\d`).ReplaceAllString(m, "*")
		return link
	}

	return ""
}

func validateChat(message *tgbotapi.Message, config *config.Config) bool {
	var chatid int64
	chatid = message.Chat.ID
	valid, _ := mw.InArray(chatid, config.Bot.ChatWhiteList)

	logrus.Infof("%v", valid, message.Chat.ID, config.Bot.ChatWhiteList)
	return valid
}
