package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	// JSONUnmarshaller is a "custom" json.Unmarshaller, to still make use of basic json unmarshalling without stackoverflows
	JSONUnmarshaller interface {
		FromJSON([]byte) error
	}

	// Interface which combnes both required interfaces
	ApiInteractor interface {
		JSONUnmarshaller
		fmt.Stringer

		Path(contestId, id string) string
		Generator() ApiInteractor
	}

	ApiTime time.Time

	ApiRelTime time.Duration

	Contest struct {
		Id         string     `json:"id"`
		Name       string     `json:"name"`
		FormalName string     `json:"formal_name"`
		StartTime  ApiTime    `json:"start_time"`
		Duration   ApiRelTime `json:"duration"`
	}

	JudgementType struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Penalty bool   `json:"penalty"`
		Solved  bool   `json:"solved"`
	}

	Problem struct {
		Id      string `json:"id"`
		Label   string `json:"label"`
		Name    string `json:"name"`
		Ordinal int    `json:"ordinal"`
	}

	Submission struct {
		Id          string     `json:"id"`
		LanguageId  string     `json:"language_id"`
		ContestTime ApiRelTime `json:"contest_time"`
		Time        ApiTime    `json:"time"`
		TeamId      string     `json:"team_id"`
		ProblemId   string     `json:"problem_id"`
		ExternalId  string     `json:"external_id"`
		EntryPoint  string     `json:"entry_point"`
	}

	Judgement struct {
		Id               string     `json:"id"`
		SubmissionId     string     `json:"submission_id"`
		JudgementTypeId  string     `json:"judgement_type_id"`
		StartContestTime ApiRelTime `json:"start_contest_time"`
		StartTime        ApiTime    `json:"start_time"`
		EndContestTime   ApiRelTime `json:"end_contest_time"`
		EndTime          ApiTime    `json:"end_time"`
	}
)

// Ensure all types adhere to required interfaces
var (
	_ ApiInteractor = new(Contest)
	_ ApiInteractor = new(JudgementType)
	_ ApiInteractor = new(Problem)
	_ ApiInteractor = new(Submission)
	_ ApiInteractor = new(Judgement)

	_ json.Unmarshaler = new(ApiTime)
	_ JSONUnmarshaller = new(ApiTime)
	_ fmt.Stringer     = new(ApiTime)

	_ json.Unmarshaler = new(ApiRelTime)
	_ JSONUnmarshaller = new(ApiRelTime)
	_ fmt.Stringer     = new(ApiRelTime)
)

func (c *Contest) FromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}

func (c Contest) String() string {
	// TODO format the starttime and duration
	return fmt.Sprintf(`
         id: %v
       name: %v
formal name: %v
 start time: %v
   duration: %v
`, c.Id, c.Name, c.FormalName, c.StartTime, c.Duration)
}

func (c Contest) Path(contestId, id string) string {
	return "contests"
}

func (c *Contest) Generator() ApiInteractor {
	return new(Contest)
}

func (jt *JudgementType) FromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, jt)
}

func (jt JudgementType) String() string {
	return fmt.Sprintf(`
         id: %v
       name: %v
    penalty: %v
     solved: %v
`, jt.Id, jt.Name, jt.Penalty, jt.Solved)
}

func (jt JudgementType) Path(contestId, id string) string {
	return "contests/" + contestId + "/judgement-types/" + id
}

func (jt JudgementType) Generator() ApiInteractor {
	return new(JudgementType)
}

func (p *Problem) FromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

func (p Problem) String() string {
	return fmt.Sprintf(`
         id: %v
      label: %v
       name: %v
    ordinal: %v
`, p.Id, p.Label, p.Name, p.Ordinal)
}

func (p Problem) Path(contestId, id string) string {
	return "contests/" + contestId + "/problems/" + id
}

func (p Problem) Generator() ApiInteractor {
	return new(Problem)
}

func (a *ApiTime) FromJSON(b []byte) error {
	return a.UnmarshalJSON(b)
}

func (a *ApiTime) UnmarshalJSON(b []byte) (err error) {
	data := strings.Trim(string(b), "\"")

	if data == "null" {
		*a = ApiTime(time.Time{})
		return
	}

	// All possible time formats we support
	var supportedTimeFormats = []string{
		// time.RFC3999 also accepts milliseconds, even though it is not officially stated
		time.RFC3339,
		// time.RFC3999 but then without the minutes of the timezone
		"2006-01-02T15:04:05Z07",
	}
	for _, supportedTimeFormat := range supportedTimeFormats {
		if t, err := time.Parse(supportedTimeFormat, data); err == nil {
			*a = ApiTime(t)
			return nil
		}
	}

	return fmt.Errorf("can not format date: %s", data)
}

func (a *ApiRelTime) FromJSON(b []byte) error {
	return a.UnmarshalJSON(b)
}

func (a *ApiRelTime) UnmarshalJSON(b []byte) (err error) {
	data := strings.Trim(string(b), "\"")
	if data == "null" {
		*a = 0
		return
	}
	re := regexp.MustCompile("(-?[0-9]{1,2}):([0-9]{2}):([0-9]{2})(.([0-9]{3}))?")
	sm := re.FindStringSubmatch(data)
	h, err := strconv.ParseInt(sm[1], 10, 64)
	if err != nil {
		return err
	}

	m, err := strconv.ParseInt(sm[2], 10, 64)
	if err != nil {
		return err
	}

	s, err := strconv.ParseInt(sm[3], 10, 64)
	if err != nil {
		return err
	}

	var ms int64 = 0
	if sm[5] != "" {
		ms, err = strconv.ParseInt(sm[5], 10, 64)
		if err != nil {
			return err
		}
	}

	*a = ApiRelTime(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second + time.Duration(ms)*time.Millisecond)

	return
}

func (s *Submission) FromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

func (s Submission) Path(contestId, id string) string {
	return "contests/" + contestId + "/submissions/" + id
}

func (s Submission) Generator() ApiInteractor {
	return new(Submission)
}

func (s Submission) String() string {
	return fmt.Sprintf(`
          id: %v
 language id: %v
contest time: %v
     team id: %v
  problem id: %v
 external id: %v
 entry_point: %v
`, s.Id, s.LanguageId, s.ContestTime, s.TeamId, s.ProblemId, s.ExternalId, s.EntryPoint)
}

func (j *Judgement) FromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, j)
}

func (j Judgement) Path(contestId, id string) string {
	return "contests/" + contestId + "/judgements/" + id
}

func (j Judgement) Generator() ApiInteractor {
	return new(Judgement)
}

func (j Judgement) String() string {
	return fmt.Sprintf(`
                id: %v
     submission id: %v
start_contest time: %v
  end_contest time: %v
 judgement_type id: %v
`, j.Id, j.SubmissionId, j.StartContestTime, j.EndContestTime, j.JudgementTypeId)
}

func (a ApiRelTime) Duration() time.Duration {
	return time.Duration(a)
}

func (a ApiTime) Time() time.Time {
	return time.Time(a)
}

func (a ApiRelTime) String() string {
	return time.Duration(a).String()
}

func (a ApiTime) String() string {
	return time.Time(a).String()
}
