# Prototype Contest API cli
This repo contains a go cli for interacting with the ICPC Contest API: https://ccs-specs.icpc.io/contest_api.

# Current Use
Currently supported commands:

## List of contests
`./contest -url <url> -c <contest-id> -u <user> -p <password> contests`

Outputs a list of the contests. Why do I have to know one contest id in
order to list all of them? No good reason, just haven't finished argument parsing yet.

## List the problems
`./contest -url <url> -c <contest-id> -u <user> -p <password> problems`

Outputs a list of the problems in this contest.

## Posting a clarification
`./contest -url <url> -c <contest-id> -u <user> -p <password> clar <text>`

Problem id is currently hardcoded to "checks" in this first pass.

## Posting a submission
`./contest -url <url> -c <contest-id> -u <user> -p <password> submit <text>`

Doesn't work yet.

## Options
- `-insecure` - skips certificate checks.


# Future CLI Design Principles

## This is (primarily) a team cli
It should be focussed on making it easy to do the things that teams need to do:
- Login & connect to a contest
- List problems
- Post a submission
- View your submissions, and their judgements
- Post a clarification
- View your clarifications, responses, and broadcasts

There are also some nice to haves - not essential, but we could add in the future:
- Current score/scoreboard
- Contest time

## Identifiers (IDs) are 'internal'
IDs can be meaningless strings up to 36 characters long. Contest ids and team ids tend to be short and useful, but that can't be expected from other types. As such, other identifying information should be used in messages as much as possible:
- Problems: label
- Submissions: problem, time
- Clarifications: problem, time

## Flags are mostly for scripting
We should try to detect things like language and problem as much as possible, and
prompt for confirmation or required choices. The flags are mostly there for
scripting (e.g. being able to do everything in one command) or edge-cases.


# Proposed Future CLI

## Flags

```
-url   base contest url
-c     contest id
-u     user
-p     password
-i     insecure (don't check certificates)
-f     force (don't prompt)
```

## Configuration


| Command | Purpose |
| ------- | ------- |
| `contest set url <url>` | Specify the base contest URL. Can also be specified via -url. |
| `contest set id <id>` | Only required in cases where there is more than one contest *and* there isn't an obvious default contest. Can also be specified via using -c. |
| `contest login <user> <password>` | Loaded from .netrc if not specified. Not required when using IP/host auto-login. Can also be specified via -u and -p. |
| `contest logout` | 'nuf said. |


## Regular Commands

| Command | Purpose |
| ------- | ------- |
| `contest list contests` | Lists all the contests. TODO - this is ugly |
| `contest list problems` | Lists all the problems in this contest. |
| `contest list clarifications` | List all clarifications this team can see: posted clarifications, responses, and broadcast messages. |
| `contest list submissions` | List all of the team's submissons and judgements. |
| `contest post-clar <problemLabel> text` | Post a clarification to the contest. |
| `contest submit [problemId] [languageId] [entry_point] file1 [<file2> <file3> ...]` | Post a submission for a problem. |


# Examples

## 'Basic' Use

```
> contest set url https://ccs/api
URL validated. 2 contests found (test, finals [current])

> contest login team47 mn3r0f
User & password valid.

> contest post-clar C "Can we assume x is never 0?"
Posting clarification to problem C [Enter to confirm]
Clarification posted successfully (clar47)
```

## Use at ICPC Finals
- Url pre-configured
- Auto-login by IP

```
> contest post-clar C "Can we assume x is never 0?"
Posting clarification to problem C [Enter to confirm]
Clarification posted successfully: clar47

> contest submit A.java
Posting solution to problem A using Java [Enter to confirm]
Submission posted successfully: 2469
```

## Tool/Scripted Use
Tools can use flags to do everything in one command:

```
> contest -url https://ccs/api -c finals -u team47 -p mn3r0f post-clar C "Can we assume x is never 0?" -f
Clarification posted successfully (clar47)
```

## Possible Future Interactive Mode
```
> contest post-clar
What problem is your clarification related to? [A/B/C/D/E/F (Enter for none)]
B
What is your clarification?
Can we assume x is never 0?
Clarification posted successfully (clar47)
```
