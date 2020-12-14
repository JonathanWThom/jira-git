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
	baseURL  = os.Getenv("ATLASSIAN_BASE_URL")
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
	for _, param := range []string{apiToken, baseURL, email} {
		if param == "" {
			return errors.New("Please set all environment variables")
		}
	}

	issueID, err := getIssueID(os.Args[1])
	if err != nil {
		return err
	}

	jc, err := newJiraClient(email, apiToken)
	if err != nil {
		return err
	}

	summary, err := jc.GetIssueSummary(issueID)
	if err != nil {
		return err
	}

	return checkoutNewBranch(issueID, summary)
}

func getIssueID(uri string) (string, error) {
	url, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}
	paths := strings.Split(url.Path, "/")

	// Probably want to handle invalid urls here
	// Something that might return 0 paths....
	// We'll get to testing soon

	return paths[len(paths)-1], nil
}

var reg = regexp.MustCompile("[^a-zA-Z0-9 ]+")

func checkoutNewBranch(id, summary string) error {
	sanitized := reg.ReplaceAllString(summary, "")
	dashed := strings.ToLower(strings.ReplaceAll(sanitized, " ", "-"))
	branchName := fmt.Sprintf("%s--%s", id, dashed)
	fmt.Println(branchName)
	// Could maybe stash here first?
	cmd := exec.Command("git", "checkout", "-b", branchName) // This could probably be replaced with something like `git-go` to avoid shelling out

	return cmd.Run()
}

type jiraClient struct {
	client *jira.Client
}

// Tell me to tell you about functional parameters. It's a nice way to add mutliple configuration, and having nice defaults
func newJiraClient(email, apiToken string) (*jiraClient, error) {
	tp := jira.BasicAuthTransport{
		Username: email,
		Password: apiToken,
	}
	client, err := jira.NewClient(tp.Client(), baseUrl)
	if err != nil {
		return nil, err
	}
	return &jiraClient{
		client: client,
	}, nil
}

func (c *jiraClient) GetIssueSummary(id string) (string, error) {
	issue, _, err := c.client.Issue.Get(id, nil)
	if err != nil {
		return "", err
	}

	return issue.Fields.Summary, nil
}
