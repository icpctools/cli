# Prototype Contest API cli
This repo contains a go cli for interacting with the ICPC Contest API: https://ccs-specs.icpc.io/contest_api.

## Current Usage
`go build .`, then

### List of contests
`./contest -url <url> -c <contest-id> -u <user> -p <password> contests`

Outputs a list of the contests. Why do I have to know one contest id in
order to list all of them? No good reason, just haven't finished argument parsing yet.

### List the problems
`./contest -url <url> -c <contest-id> -u <user> -p <password> problems`

Outputs a list of the problems in this contest.

### Posting a clarification
`./contest -url <url> -c <contest-id> -u <user> -p <password> clar <text>`

Problem id is currently hardcoded to "checks" in this first pass.

### Posting a submission
`./contest -url <url> -c <contest-id> -u <user> -p <password> submit <text>`

Doesn't work yet.

### Options
- `-insecure` - skips certificate checks.

# Proposed Future CLI

## Commands

### Basic Commands

`contest url <url>`

Specify the base contest URL. Can also be specified directly in other commands by using -url.

`contest login <user> <password>`

Loaded from .netrc if not specified. Not required when using IP/host auto-login. Can also be specified directly in other commands by using -u and -p.

`contest id <id>`

Only required in cases where there is more than one contest *and* there isn't an obvious default contest. Can also be specified directly in other commands by using -c.
e.g. at finals there might be three contests, but we'll automatically pick the correct current one based on starting time/state.

### Listing the contests
TBD, would list all the contests with enough details that you could know which contest to use for `contest id`. `contest list` or `contest list contests`?

### Listing the problems

TBD. `contest problems` or `contest list problems`?

### Posting a clarification
`contest clar <problemLabel> <text>`

Using all optional params

`contest -url <url> -c <contestId> -u <user> -p <password> clar <problemId> <text>`

### Listing clarifications
TBD, would show all submitted clarifications, responses, and broadcast messages. `contest list clars`?

### Posting a submission
`contest submit -p <problemId> -l <languageId> -e <entry_point> <file1> [<file2> <file3> ...]`

### Listing submissions
TBD, would show all team's submissons and judgements. `contest list submissions`?

## Examples

### 'Basic' Use
`contest url https://ccs/api`

> URL validated. 2 contests found (test, finals)

`contest login team47 mn3r0f`

> User & password correct

`contest clar C "Can we assume x is never 0?"`

> Clarification posted successfully: clar47

### Standalone Use
`contest -url https://ccs/api -c finals -u team47 -p mn3r0f clar C "Can we assume x is never 0?"`

> Clarification posted successfully: clar47

### Use at ICPC Finals
- Url pre-configured
- Auto-login by IP

`contest clar C "Can we assume x is never 0?"`

> Clarification posted successfully: clar47

### Interactive Mode
`contest clar`

> 4 Problems found (A, B, C, D), pick one: `B`

> What is your clarification?: `Can we assume x is never 0?`

> Clarification posted successfully: clar47