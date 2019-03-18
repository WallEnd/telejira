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
	Project     Project   `json:"project"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	IssueType   IssueType `json:"issuetype"`
	Status      Status    `json:"status"`
}

type Project struct {
	Key string `json:"key"`
}

type Status struct {
	Name string `json:"name"`
}

type IssueType struct {
	Name string `json:"name"`
}

type JiraIssue struct {
	Fields Fields `json:"fields"`
	Key    string `json:"key"`
}

type NewIssue struct {
	Fields Fields `json:"fields"`
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

	return "*[" + jiraissue.Fields.IssueType.Name + "]* [" + jiraissue.Key + "](" + config.JiraUrl + "/browse/" + jiraissue.Key + ") - *" + jiraissue.Fields.Summary + "*" +
		"\n *(status:" + jiraissue.Fields.Status.Name + ")* \n" + jiraissue.Fields.Description

}

func Init(data Config) {
	config = data
}
