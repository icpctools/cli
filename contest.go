package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

// Current use:
//   contest -url <baseURL> -u <user> -p <password> contests
//   contest -url <baseURL> -c contestId -u <user> -p <password> problems
//   contest -url <baseURL> -c contestId -u <user> -p <password> clar <text>
//   contest -url <baseURL> -c contestId -u <user> -p <password> submit <text>
// -insecure is optional, the rest isn't. Problem id is hardcoded to "checks" for now

type (
	Contest struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	Problem struct {
		Label string `json:"label"`
		Name  string `json:"name"`
	}

	// Implementation of the http.RoundTripper interface, used for always adding basic-auth
	basicAuthTransport struct {
		T http.RoundTripper
	}
)

// RoundTrip adds the basic auth headers
func (b basicAuthTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}

	return b.T.RoundTrip(request)
}

var (
	baseURL   string
	contestId string
	user      string
	password  string

	insecure bool

	// Ensure basicAuthTransport adheres to the interface
	_ http.RoundTripper = basicAuthTransport{}

	errUnauthorized = errors.New("request not authorized")
	errNotFound     = errors.New("object not found")
)

func init() {
	flag.StringVar(&baseURL, "url", "", "the base 'Contest URL'")
	flag.StringVar(&contestId, "c", "", "the 'contest id'")
	flag.StringVar(&user, "u", "", "the 'user'")
	flag.StringVar(&password, "p", "", "the 'password'")
	insecure = *flag.Bool("insecure", false, "whether the connection is secure")
	flag.Parse()

	// Strip trailing slash of the baseurl
	baseURL = strings.TrimSuffix(baseURL, "/")
}

func main() {
	if len(flag.Args()) == 0 {
		fmt.Println("No command given")
		return
	}

	if baseURL == "" {
		fmt.Println("No base contest URL")
		return
	}

	if contestId == "" {
		fmt.Println("No contest id")
		return
	}

	// Create a transport for insecure communication and adding of basic-auth headers
	transport := http.DefaultTransport
	if insecure {
		fmt.Println("Insecure")
		transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Create client and transport for basic auth
	client := http.Client{
		Transport: basicAuthTransport{transport},
	}

	// flag.Arg always return empty string if arg is not set
	var text = flag.Arg(1)
	var hasText = text != ""

	// Switch is preferred over large else-if chains, though sometimes (as with this nested case) else-if is clearly more
	// legible
	switch flag.Arg(0) {
	case "clar":
		if !hasText {
			fmt.Println("Clarification text not provided correctly")
		} else if clarificationId, err := postClarification(client, "checks", text); err != nil {
			fmt.Println("Error submitting clarification: ", err)
		} else {
			fmt.Println("Clarification submitted successfully! ", clarificationId)
		}
	case "submit":
		if !hasText {
			fmt.Println("Submission not provided correctly")
		} else if submissionId, err := postSubmission(client, text, text, text); err != nil {
			fmt.Println("Error submitting:", err)
		} else {
			fmt.Println("Submitted successfully!", submissionId)
		}
	case "contests":
		if contests, err := getContests(client); err != nil {
			fmt.Println("Error getting contests:", err)
		} else {
			fmt.Println("Contests found successfully!", contests)
			fmt.Printf("Contests: %+v", contests)
		}
	case "problems":
		if problems, err := getProblems(client); err != nil {
			fmt.Println("Error getting problems:", err)
		} else {
			fmt.Println("Problems found successfully!", problems)
			fmt.Printf("Problems: %+v", problems)
		}
	default:
		fmt.Println("Only 'clar', 'submit', 'contests', and 'problems' commands are supported")
	}
}

func getContests(client http.Client) ([]Contest, error) {
	var contests []Contest

	resp, err := client.Get(baseURL + "/contests")
	// Body is always non-nil, ensure it will always be closed
	defer resp.Body.Close()
	if err != nil {
		return contests, err
	} else if err := handleStatus(resp.StatusCode); err != nil {
		return contests, err
	}

	// Parse response, no need for specific control flow on error
	return contests, json.NewDecoder(resp.Body).Decode(&contests)
}

func getProblems(client http.Client) ([]Problem, error) {
	var problems []Problem

	resp, err := client.Get(baseURL + "/contests/" + contestId + "/problems")
	// Body is always non-nil, ensure it will always be closed
	defer resp.Body.Close()
	if err != nil {
		return problems, err
	} else if err := handleStatus(resp.StatusCode); err != nil {
		return problems, err
	}

	// Parse response, no need for specific control flow on error
	return problems, json.NewDecoder(resp.Body).Decode(&problems)
}

func postClarification(client http.Client, problemId string, text string) (string, error) {
	var clarificationId string

	var buf = new(bytes.Buffer)
	enc := json.NewEncoder(buf)

	// Perhaps replace this with a struct?
	err := enc.Encode(map[string]string{
		"problem_id": problemId,
		"text":       text,
	})

	if err != nil {
		return clarificationId, err
	}

	resp, err := client.Post(baseURL+"/contests/"+contestId+"/clarifications", "application/json", buf)
	defer resp.Body.Close()
	if err != nil {
		return clarificationId, err
	} else if err := handleStatus(resp.StatusCode); err != nil {
		return clarificationId, err
	}

	// parse response
	return clarificationId, json.NewDecoder(resp.Body).Decode(&clarificationId)
}

func postSubmission(client http.Client, problemId string, languageId string, file string) (string, error) {
	var submissionId string

	var buf = new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(map[string]string{
		"problem_id":  problemId,
		"language_id": languageId,
		"file":        file,
	})

	if err != nil {
		return submissionId, err
	}

	resp, err := client.Post(baseURL+"/contests/"+contestId+"/submissions", "application/json", buf)
	if err != nil {
		return submissionId, err
	} else if err := handleStatus(resp.StatusCode); err != nil {
		return "", err
	}

	// parse response
	return submissionId, json.NewDecoder(resp.Body).Decode(&submissionId)
}

func handleStatus(status int) error {
	switch status {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return errUnauthorized
	case http.StatusNotFound:
		return errNotFound
	default:
		return fmt.Errorf("invalid statuscode received: %d", status)
	}
}
