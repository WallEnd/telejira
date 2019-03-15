package jiraapi

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	JiraLogin           string
	JiraPass            string
	JiraUrl             string
	JiraApiKeySearchUrl string
}

type Fields struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type JiraIssue struct {
	Fields Fields `json:"fields"`
	Key    string `json:"key"`
}

var config Config

func GetIssueLink(key string) string {
	var jiraissue JiraIssue
	req, err := http.NewRequest("GET", config.JiraApiKeySearchUrl+key, nil)
	req.SetBasicAuth(config.JiraLogin, config.JiraPass)

	client := &http.Client{}

	if err != nil {
		logrus.Fatal(err)
	}

	resp, resperr := client.Do(req)

	if resperr != nil {
		logrus.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	_ = json.Unmarshal(b, &jiraissue)

	return "<a href=\"" + config.JiraUrl + "/browse/" + jiraissue.Key + "\">" + jiraissue.Key + "</a> - <b>" + jiraissue.Fields.Summary + "</b>" +
		"\n " + jiraissue.Fields.Description

}

func Init(data Config) {
	config = data
}
