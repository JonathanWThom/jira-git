# Jira-Git

Open an idiomatically named branch using a Jira Issue URL.

### Usage

```
$ cd directory-of-your-project
$ jira-git https://your-jira-instance/browse/TEAM-1234
```

This will run `git checkout -b TEAM-1234--summary-of-your-ticket`.

### Installation

1. Install Go
2. Get Atlassian API Token from
   [here](https://id.atlassian.com/manage-profile/security/api-tokens).
3. Set shell environment variables:
  ```
  export ATLASSIAN_API_TOKEN=token-from-above
  export ATLASSIAN_BASE_URL=https://your-jira-instance.atlassian.net
  export ATLASSIAN_EMAIL=email-you-use-to-login@domain.com
  ```
4. Clone this, `cd` in, and `go install`. Ideally the `jira-git` command will
   not be somewhere in your `PATH`.

### License

MIT

