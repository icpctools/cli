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

// Implementation of the http.RoundTripper interface, used for always adding basic-auth
type basicAuthTransport struct {
	T http.RoundTripper
}

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

	interactors = map[string]ApiInteractor{
		"contests":        new(Contest),
		"problems":        new(Problem),
		"submissions":     new(Submission),
		"judgement-types": new(JudgementType),
		"judgements":      new(Judgement),
	}
)

func init() {
	flag.StringVar(&baseURL, "url", "", "the base 'Contest URL'")
	flag.StringVar(&contestId, "c", "", "the 'contest id'")
	flag.StringVar(&user, "u", "", "the 'user'")
	flag.StringVar(&password, "p", "", "the 'password'")
	flag.BoolVar(&insecure, "insecure", false, "whether the connection is secure")
}

func main() {
	flag.Parse()

	// Strip trailing slash of the baseurl
	baseURL = strings.TrimSuffix(baseURL, "/")

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
	id := flag.Arg(1)
	switch typ := flag.Arg(0); typ {
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
	case "contests", "problems", "submissions", "judgement-types", "judgements":
		title := strings.Title(typ)
		if objects, err := getObjects(client, interactors[typ], id); err != nil {
			fmt.Printf("Error getting %v: %v\n", title, err)
		} else {
			fmt.Printf("%v (%d):", title, len(objects))
			for _, o := range objects {
				fmt.Print(o)
			}
		}
	default:
		fmt.Println("Only 'clar', 'submit', 'contests', 'submissions', and 'problems' commands are supported")
	}
}

func getObjects(client http.Client, interactor ApiInteractor, id string) ([]ApiInteractor, error) {
	resp, err := client.Get(baseURL + "/" + interactor.Path(contestId, id))
	if err != nil {
		return nil, err
	}

	// Body is not-nil, ensure it will always be closed
	defer resp.Body.Close()

	if err := handleStatus(resp.StatusCode); err != nil {
		return nil, err
	}

	// Some json should be returned, construct a decoder
	decoder := json.NewDecoder(resp.Body)

	// If id is not empty, only a single instance is expected to be returned
	if id != "" {
		in := interactor.Generator()
		return []ApiInteractor{in}, decoder.Decode(in)
	}

	// We read everything into a slice of
	var temp []json.RawMessage
	if err := decoder.Decode(&temp); err != nil {
		return nil, err
	}

	// Create the actual slice to return
	ret := make([]ApiInteractor, len(temp))
	for k, v := range temp {
		// Generate a new interactor
		in := interactor.Generator()
		if err := in.FromJSON(v); err != nil {
			return ret, err
		}

		ret[k] = in
	}

	return ret, nil
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
