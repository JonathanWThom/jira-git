package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	jira "github.com/andygrunwald/go-jira"
)

var (
	apiToken = os.Getenv("ATLASSIAN_API_TOKEN")
	baseUrl  = os.Getenv("ATLASSIAN_BASE_URL")
	email    = os.Getenv("ATLASSIAN_EMAIL")
)

func main() {
	if err := Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Run(args []string) error {
	if len(args) != 2 {
		return errors.New("Please pass Issue URL as argument")
	}
	issueID, err := getIssueID(os.Args[1])
	if err != nil {
		return err
	}

	summary, err := getIssueSummary(issueID)
	if err != nil {
		return err
	}

	return checkoutNewBranch(issueID, summary)
}

func getIssueID(uri string) (string, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	paths := strings.Split(url.Path, "/")

	return paths[len(paths)-1], nil
}

func checkoutNewBranch(id, summary string) error {
	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		return err
	}
	sanitized := reg.ReplaceAllString(summary, "")
	dashed := strings.ToLower(strings.ReplaceAll(sanitized, " ", "-"))
	branchName := fmt.Sprintf("%s--%s", id, dashed)
	fmt.Println(branchName)
	// Could maybe stash here first?
	cmd := exec.Command("git", "checkout", "-b", branchName)

	return cmd.Run()
}

func getIssueSummary(id string) (string, error) {
	for _, param := range []string{apiToken, baseUrl, email} {
		if param == "" {
			return "", errors.New("Please set all environment variables")
		}
	}

	tp := jira.BasicAuthTransport{
		Username: email,
		Password: apiToken,
	}

	client, err := jira.NewClient(tp.Client(), baseUrl)
	if err != nil {
		return "", err
	}
	issue, _, err := client.Issue.Get(id, nil)
	if err != nil {
		return "", err
	}

	return issue.Fields.Summary, nil
}
