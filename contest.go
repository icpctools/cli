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
//   contest -url <baseURL> -u <user> -p <password> contests
//   contest -url <baseURL> -c contestId -u <user> -p <password> problems
//   contest -url <baseURL> -c contestId -u <user> -p <password> clar <text>
//   contest -url <baseURL> -c contestId -u <user> -p <password> submit <text>
// -insecure is optional, the rest isn't. Problem id is hardcoded to "checks" for now

var baseURL string
var contestId string
var user string
var password string

type Contest struct {
	Id   string
	Name string
}

type Problem struct {
	Label string
	Name  string
}

func getJson(resp http.Response, target interface{}) error {
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func getContests(client http.Client) ([]Contest, error) {
	req, err := http.NewRequest("GET", baseURL+"/contests/", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// parse response
	var contests []Contest
	getJson(*resp, &contests)
	return contests, nil
}

func getProblems(client http.Client) ([]Problem, error) {
	req, err := http.NewRequest("GET", baseURL+"/contests/"+contestId+"/problems", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// parse response
	var problems []Problem
	getJson(*resp, &problems)
	return problems, nil
}

func postClarification(client http.Client, problemId string, text string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"problem_id": problemId,
		"text":       text,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/contests/"+contestId+"/clarifications", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	// parse response
	var clarificationId string
	getJson(*resp, &clarificationId)
	if err != nil {
		return "", err
	}
	return clarificationId, nil
}

func postSubmission(client http.Client, problemId string, languageId string, file string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"problem_id":  problemId,
		"language_id": languageId,
		"file":        file,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/contests/"+contestId+"/submissions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	var submissionId string
	getJson(*resp, &submissionId)
	if err != nil {
		return "", err
	}
	return submissionId, nil
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.SetBasicAuth("team1", "team11")
	return nil
}

func main() {
	flag.StringVar(&baseURL, "url", "", "the base 'Contest URL'") // TODO if there is a trailing slash we should remove it
	flag.StringVar(&contestId, "c", "", "the 'contest id'")
	flag.StringVar(&user, "u", "", "the 'user'")
	flag.StringVar(&password, "p", "", "the 'password'")
	insecure := flag.Bool("insecure", false, "whether the connection is secure")
	flag.Parse()

	if baseURL == "" {
		fmt.Println("No base contest URL")
		return
	}

	if contestId == "" {
		fmt.Println("No contest id")
		return
	}

	if len(flag.Args()) == 0 {
		fmt.Println("No command given")
		return
	}

	if *insecure {
		fmt.Println("Insecure")
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{
		//Jar: cookieJar,
		CheckRedirect: redirectPolicyFunc,
	}

	if flag.Args()[0] == "clar" {
		if len(flag.Args()) != 2 {
			fmt.Println("Clarification text not provided correctly")
			return
		}

		text := flag.Args()[1]
		clarificationId, err := postClarification(*client, "checks", text)
		if err != nil {
			fmt.Println("Error submitting clarification: ", err)
			return
		} else {
			fmt.Println("Clarification submitted successfully! ", clarificationId)
		}
	} else if flag.Args()[0] == "submit" {
		if len(flag.Args()) != 2 {
			fmt.Println("Submission not provided correctly")
			return
		}

		text := flag.Args()[1]
		submissionId, err := postSubmission(*client, text, text, text)
		if err != nil {
			fmt.Println("Error submitting: ", err)
			return
		} else {
			fmt.Println("Submitted successfully! ", submissionId)
		}
	} else if flag.Args()[0] == "contests" {
		contests, err := getContests(*client)
		if err != nil {
			fmt.Println("Error getting contests: ", err)
			return
		} else {
			fmt.Println("Contests found successfully! ", contests)
			fmt.Printf("Contests: %+v", contests)
		}
	} else if flag.Args()[0] == "problems" {
		problems, err := getProblems(*client)
		if err != nil {
			fmt.Println("Error getting problems: ", err)
			return
		} else {
			fmt.Println("Problems found successfully! ", problems)
			fmt.Printf("Problems: %+v", problems)
		}
	} else {
		fmt.Println("Only 'clar' command is supported")
	}
}
