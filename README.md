# Prototype Contest API cli
This repo contains a go cli for interacting with the ICPC Contest API: https://ccs-specs.icpc.io/contest_api.

## Usage
`go build .`, then

`./contest -c <contestURL> -u <user> -p <password> -insecure clar <text>`

`-insecure` is optional, the rest isn't. Problem id is hardcoded to "checks" in this first pass.

## Possible Future CLI
`contest set <url>`

`contest login <user> <password>`

`contest clar -p <problemId> <text>`

`contest submit -p <problemId> -l <languageId> -e <entry_point> <file>`

Interactive mode?