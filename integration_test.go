// +build integration

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
)

// Fake response that looks like that of JIRA - used the roundtripper to get an example response
const fakeResponse = `{
	"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations,customfield_10212.properties",
	"id":"35159",
	"self":"https://example.atlassian.net/rest/api/2/issue/35159",
	"key":"TOOL-936",
	"fields":{
		 "statuscategorychangedate":"2020-12-08T03:46:14.032-0800",
		 "customfield_10190":null,
		 "customfield_10191":null,
		 "customfield_10192":null,
		 "customfield_10193":null,
		 "customfield_10194":null,
		 "customfield_10195":null,
		 "customfield_10196":null,
		 "customfield_10197":null,
		 "fixVersions":[

		 ],
		 "customfield_10198":null,
		 "customfield_10199":null,
		 "resolution":null,
		 "customfield_10104":null,
		 "customfield_10105":null,
		 "customfield_10107":null,
		 "lastViewed":"2020-12-14T00:12:17.591-0800",
		 "customfield_10180":null,
		 "customfield_10181":null,
		 "customfield_10182":null,
		 "customfield_10183":null,
		 "customfield_10184":null,
		 "customfield_10185":null,
		 "customfield_10186":null,
		 "customfield_10187":null,
		 "priority":{
				"self":"https://example.atlassian.net/rest/api/2/priority/10000",
				"iconUrl":"https://example.atlassian.net/secure/viewavatar?size=medium&avatarId=10671&avatarType=issuetype",
				"name":"N/A",
				"id":"10000"
		 },
		 "customfield_10188":null,
		 "customfield_10189":null,
		 "labels":[

		 ],
		 "customfield_10103":null,
		 "timeestimate":null,
		 "aggregatetimeoriginalestimate":null,
		 "versions":[

		 ],
		 "issuelinks":[

		 ],
		 "assignee":{
				"self":"https://example.atlassian.net/rest/api/2/user?accountId=accountID",
				"accountId":"accountID",
				"emailAddress":"max@example.com",
				"avatarUrls":{
					 "48x48":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png",
					 "24x24":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png",
					 "16x16":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png",
					 "32x32":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png"
				},
				"displayName":"Max Man",
				"active":true,
				"timeZone":"America/Los_Angeles",
				"accountType":"atlassian"
		 },
		 "status":{
				"self":"https://example.atlassian.net/rest/api/2/status/10313",
				"description":"This status is managed internally by Jira Software",
				"iconUrl":"https://example.atlassian.net/",
				"name":"Code Review",
				"id":"10313",
				"statusCategory":{
					 "self":"https://example.atlassian.net/rest/api/2/statuscategory/4",
					 "id":4,
					 "key":"indeterminate",
					 "colorName":"yellow",
					 "name":"In Progress"
				}
		 },
		 "components":[

		 ],
		 "customfield_10170":null,
		 "customfield_10171":null,
		 "customfield_10172":null,
		 "customfield_10173":null,
		 "customfield_10174":null,
		 "customfield_10175":null,
		 "customfield_10177":null,
		 "customfield_10210":null,
		 "customfield_10211":null,
		 "customfield_10178":null,
		 "customfield_10179":null,
		 "customfield_10212":null,
		 "customfield_10213":null,
		 "customfield_10205":null,
		 "aggregatetimeestimate":null,
		 "customfield_10209":null,
		 "creator":{
				"self":"https://example.atlassian.net/rest/api/2/user?accountId=accountID",
				"accountId":"accountID",
				"emailAddress":"max@example.com",
				"avatarUrls":{
					 "48x48":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png",
					 "24x24":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png",
					 "16x16":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png",
					 "32x32":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png"
				},
				"displayName":"Max Man",
				"active":true,
				"timeZone":"America/Los_Angeles",
				"accountType":"atlassian"
		 },
		 "subtasks":[

		 ],
		 "customfield_10160":null,
		 "customfield_10161":null,
		 "customfield_10162":null,
		 "customfield_10163":null,
		 "reporter":{
				"self":"https://example.atlassian.net/rest/api/2/user?accountId=accountID",
				"accountId":"accountID",
				"emailAddress":"max@example.com",
				"avatarUrls":{
					 "48x48":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png",
					 "24x24":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png",
					 "16x16":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png",
					 "32x32":"https://secure.gravatar.com/avatar/uniqID?d=https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/initials/MM-5.png"
				},
				"displayName":"Max Man",
				"active":true,
				"timeZone":"America/Los_Angeles",
				"accountType":"atlassian"
		 },
		 "customfield_10164":null,
		 "aggregateprogress":{
				"progress":0,
				"total":0
		 },
		 "customfield_10165":null,
		 "customfield_10045":null,
		 "customfield_10166":null,
		 "customfield_10167":null,
		 "customfield_10200":null,
		 "customfield_10168":null,
		 "customfield_10159":null,
		 "progress":{
				"progress":0,
				"total":0
		 },
		 "votes":{
				"self":"https://example.atlassian.net/rest/api/2/issue/TOOL-936/votes",
				"votes":0,
				"hasVoted":false
		 },
		 "worklog":{
				"startAt":0,
				"maxResults":20,
				"total":0,
				"worklogs":[

				]
		 },
		 "issuetype":{
				"self":"https://example.atlassian.net/rest/api/2/issuetype/10002",
				"id":"10002",
				"description":"A task that needs to be done.",
				"iconUrl":"https://example.atlassian.net/secure/viewavatar?size=medium&avatarId=10318&avatarType=issuetype",
				"name":"Task",
				"subtask":false,
				"avatarId":10318
		 },
		 "timespent":null,
		 "customfield_10151":null,
		 "customfield_10152":null,
		 "project":{
				"self":"https://example.atlassian.net/rest/api/2/project/10129",
				"id":"10129",
				"key":"TOOL",
				"name":"Internal Tools",
				"projectTypeKey":"software",
				"simplified":false,
				"avatarUrls":{
					 "48x48":"https://example.atlassian.net/secure/projectavatar?pid=10129&avatarId=10405",
					 "24x24":"https://example.atlassian.net/secure/projectavatar?size=small&s=small&pid=10129&avatarId=10405",
					 "16x16":"https://example.atlassian.net/secure/projectavatar?size=xsmall&s=xsmall&pid=10129&avatarId=10405",
					 "32x32":"https://example.atlassian.net/secure/projectavatar?size=medium&s=medium&pid=10129&avatarId=10405"
				},
				"projectCategory":{
					 "self":"https://example.atlassian.net/rest/api/2/projectCategory/10000",
					 "id":"10000",
					 "description":"Amazing Description",
					 "name":"Amazing Name"
				}
		 },
		 "customfield_10153":null,
		 "customfield_10154":null,
		 "customfield_10155":null,
		 "aggregatetimespent":null,
		 "customfield_10156":null,
		 "customfield_10158":null,
		 "customfield_10027":null,
		 "customfield_10148":null,
		 "customfield_10149":{
				"self":"https://example.atlassian.net/rest/api/2/customFieldOption/10449",
				"value":"No",
				"id":"10449"
		 },
		 "resolutiondate":null,
		 "workratio":-1,
		 "issuerestriction":{
				"issuerestrictions":{

				},
				"shouldDisplay":false
		 },
		 "watches":{
				"self":"https://example.atlassian.net/rest/api/2/issue/TOOL-936/watchers",
				"watchCount":1,
				"isWatching":true
		 },
		 "created":"2020-12-08T03:46:06.988-0800",
		 "customfield_10140":null,
		 "customfield_10141":null,
		 "customfield_10020":null,
		 "customfield_10142":null,
		 "customfield_10021":null,
		 "customfield_10022":"0|i03psf:",
		 "customfield_10143":null,
		 "customfield_10144":null,
		 "customfield_10023":[

		 ],
		 "customfield_10145":null,
		 "customfield_10146":null,
		 "customfield_10026":null,
		 "customfield_10147":null,
		 "customfield_10137":null,
		 "customfield_10016":null,
		 "customfield_10017":null,
		 "customfield_10138":null,
		 "customfield_10018":{
				"hasEpicLinkFieldDependency":false,
				"showField":false,
				"nonEditableReason":{
					 "reason":"EPIC_LINK_SHOULD_BE_USED",
					 "message":"To set an epic as the parent, use the epic link instead"
				}
		 },
		 "customfield_10019":null,
		 "updated":"2020-12-08T03:46:16.607-0800",
		 "timeoriginalestimate":null,
		 "customfield_10130":null,
		 "description":"Cool description",
		 "customfield_10010":null,
		 "customfield_10134":null,
		 "customfield_10135":null,
		 "customfield_10014":"TOOL-831",
		 "customfield_10015":null,
		 "customfield_10136":null,
		 "timetracking":{

		 },
		 "customfield_10005":null,
		 "customfield_10126":null,
		 "customfield_10127":null,
		 "customfield_10006":null,
		 "customfield_10128":null,
		 "security":null,
		 "customfield_10007":null,
		 "customfield_10008":null,
		 "customfield_10129":null,
		 "attachment":[

		 ],
		 "customfield_10009":null,
		 "summary":"Test Branch",
		 "customfield_10120":null,
		 "customfield_10000":"{pullrequest={dataType=pullrequest, state=MERGED, stateCount=1}, json={\"cachedValue\":{\"errors\":[],\"summary\":{\"pullrequest\":{\"overall\":{\"count\":1,\"lastUpdated\":\"2020-12-08T06:04:27.000-0800\",\"stateCount\":1,\"state\":\"MERGED\",\"dataType\":\"pullrequest\",\"open\":false},\"byInstanceType\":{\"GitHub\":{\"count\":1,\"name\":\"GitHub\"}}}}},\"isStale\":true}}",
		 "customfield_10121":null,
		 "customfield_10122":null,
		 "customfield_10001":null,
		 "customfield_10002":null,
		 "customfield_10123":null,
		 "customfield_10003":null,
		 "customfield_10124":null,
		 "customfield_10004":null,
		 "customfield_10125":null,
		 "environment":null,
		 "duedate":null,
		 "comment":{
				"comments":[

				],
				"maxResults":0,
				"total":0,
				"startAt":0
		 }
	}
}`

func TestIntegration(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, fakeResponse)
		}),
	)
	defer srv.Close()
	apiToken = "token"
	baseURL = srv.URL
	email = "email@example.com"

	err := Run([]string{"main", "http://www.google.com/foo"})
	if err != nil {
		t.Fatal(err)
	}
	var buff bytes.Buffer
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Stdout = &buff
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	expected := "foo--test-branch"
	if buff.String() != expected+"\n" {
		t.Fatalf("expceted %v, got %v", expected, buff.String())
	}

	defer func() {
		cmd := exec.Command("git", "checkout", "-")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		rmCmd := exec.Command("git", "branch", "-D", expected)
		if err := rmCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}()
}
