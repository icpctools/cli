package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
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
		Id         string `json:"id"`
		Name       string `json:"name"`
		FormalName string `json:"formal_name"`
		StartTime  string `json:"start_time"`
		Duration   string `json:"duration"`
	}

	Problem struct {
		Id      string `json:"id"`
		Label   string `json:"label"`
		Name    string `json:"name"`
		Ordinal int    `json:"ordinal"`
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

// parse a relative time (regex: "(-)?(h)*h:mm:ss(.uuu)?") and convert to integer milliseconds
func parseRelTime(relTime string) int64 {
	re := regexp.MustCompile("(-?[0-9]{1,2}):([0-9]{2}):([0-9]{2})(.[0-9]{3})?")
	sm := re.FindStringSubmatch(relTime)
	h, err := strconv.ParseInt(sm[1], 10, 64)
	if err != nil {
		return 0
	}
	m, err := strconv.ParseInt(sm[2], 10, 64)
	if err != nil {
		return 0
	}

	s, err := strconv.ParseInt(sm[3], 10, 64)
	if err != nil {
		return 0
	}

	return s*1000 + m*60000 + h*3600000
}

// format a relative time in milliseconds into a human-readable string
func formatRelTime(relTime int64) string {
	h := relTime / 3600000
	relTime -= h * 3600000
	m := relTime / 60000
	relTime -= m * 60000
	s := relTime / 1000

	var time = ""
	if h > 0 {
		time += strconv.FormatInt(h, 10) + "h"
	}
	if m > 0 {
		time += strconv.FormatInt(m, 10) + "m"
	}
	if s > 0 {
		time += strconv.FormatInt(s, 10) + "s"
	}
	return time
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
	flag.BoolVar(&insecure, "insecure", false, "whether the connection is secure")
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
			fmt.Printf("Contests (%d):\n", len(contests))
			for i := range contests {
				c := contests[i]
				fmt.Printf("  %s: %s\n", c.Id, c.Name)
				dur := parseRelTime(c.Duration)
				fmt.Printf("     %s starting at %s\n", formatRelTime(dur), c.StartTime)
			}
		}
	case "problems":
		if problems, err := getProblems(client); err != nil {
			fmt.Println("Error getting problems:", err)
		} else {
			fmt.Printf("Problems (%d):\n", len(problems))
			for i := range problems {
				p := problems[i]
				fmt.Printf("  %s: %s\n", p.Label, p.Name)
			}
		}
	default:
		fmt.Println("Only 'clar', 'submit', 'contests', and 'problems' commands are supported")
	}
}

func getContests(client http.Client) ([]Contest, error) {
	var contests []Contest

	resp, err := client.Get(baseURL + "/contests")
	// Body is always non-nil, ensure it will always be closed
	if err != nil {
		return contests, err
	}

	// Body is not-nil, ensure it will always be closed
	defer resp.Body.Close()

	if err := handleStatus(resp.StatusCode); err != nil {
		return contests, err
	}

	// Parse response, no need for specific control flow on error
	return contests, json.NewDecoder(resp.Body).Decode(&contests)
}

func getProblems(client http.Client) ([]Problem, error) {
	var problems []Problem

	resp, err := client.Get(baseURL + "/contests/" + contestId + "/problems")
	// Body is always non-nil, ensure it will always be closed
	if err != nil {
		return problems, err
	}

	// Body is not-nil, ensure it will always be closed
	defer resp.Body.Close()

	if err := handleStatus(resp.StatusCode); err != nil {
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
	if err != nil {
		return clarificationId, err
	}

	// Body is not-nil, ensure it will always be closed
	defer resp.Body.Close()

	if err := handleStatus(resp.StatusCode); err != nil {
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
	}

	// Body is not-nil, ensure it will always be closed
	defer resp.Body.Close()

	if err := handleStatus(resp.StatusCode); err != nil {
		return submissionId, err
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
