package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

// Current use:
//   contest -c <contestURL> -u <user> -p <password> -insecure clar <text>
// -insecure is optional, the rest isn't. Problem id is hardcoded to "checks" for now

var contestURL string
var user string
var password string

func postClarification(client http.Client, problemId string, text string) (*http.Response, error) {
	requestBody, err := json.Marshal(map[string]string{
		"problem_id": problemId,
		"text":       text,
	})
	if err != nil {
		fmt.Println("Error 1: ", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", contestURL+"/clarifications", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error 2: ", err)
		return nil, err
	}
	req.SetBasicAuth(user, password)
	return client.Do(req)
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.SetBasicAuth("team1", "team11")
	return nil
}

func main() {
	flag.StringVar(&contestURL, "c", "", "the 'Contest URL'")
	flag.StringVar(&user, "u", "", "the 'user'")
	flag.StringVar(&password, "p", "", "the 'password'")
	insecure := flag.Bool("insecure", false, "whether the connection is secure")
	flag.Parse()

	if contestURL == "" {
		fmt.Println("No contest URL")
		return
	}

	if len(flag.Args()) == 0 {
		fmt.Println("No command given")
		return
	}

	if flag.Args()[0] != "clar" {
		fmt.Println("Only 'clar' command is supported")
		return
	}

	if len(flag.Args()) != 2 {
		fmt.Println("Clarification text not provided correctly")
		return
	}
	text := flag.Args()[1]

	if *insecure {
		fmt.Println("Insecure")
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{
		//Jar: cookieJar,
		CheckRedirect: redirectPolicyFunc,
	}

	resp, err := postClarification(*client, "checks", text)

	if err != nil {
		fmt.Println("Error 3: ", err)
		return
	} else {
		if resp.StatusCode == 200 {
			fmt.Println("Clarification submitted successfully!")
		} else {
			fmt.Println("Error submitting clarification: ", resp.StatusCode)
		}
	}
}
