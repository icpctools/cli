package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var submissionsCommand = &cobra.Command{
	Use:     "submissions",
	Short:   "Shows past submissions and their judgements",
	RunE:    submissions,
	PreRunE: configHelper("baseurl"),
}

func submissions(cmd *cobra.Command, args []string) error {
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	// Get the problems, languages, and judgementTypes
	problems, err := api.Problems()
	if err != nil {
		return fmt.Errorf("could not get problems; %w", err)
	}

	languages, err := api.Languages()
	if err != nil {
		return fmt.Errorf("could not get languages; %w", err)
	}

	judgementTypes, err := api.JudgementTypes()
	if err != nil {
		return fmt.Errorf("could not get judgement types; %w", err)
	}

	submissions, err := api.Submissions()
	if err != nil {
		return fmt.Errorf("could not get submissions; %w", err)
	}

	judgements, err := api.Judgements()
	if err != nil {
		return fmt.Errorf("could not get judgements; %w", err)
	}

	teamId := "1"

	// sort by submission time
	sort.Slice(submissions, func(i, j int) bool {
		return submissions[i].ContestTime.Duration() < submissions[j].ContestTime.Duration()
	})

	count := 0
	for _, s := range submissions {
		if strings.EqualFold(s.TeamId, teamId) {
			count++
		}
	}

	fmt.Printf("\nSubmissions (%d):\n", count)
	var table = Table{}
	table.Header = []string{"Time", "Problem", "Language", "Judgement Time", "Judgement"}
	table.Align = []int{ALIGN_RIGHT, ALIGN_LEFT, ALIGN_LEFT, ALIGN_RIGHT, ALIGN_LEFT}
	for _, s := range submissions {
		if s.TeamId == teamId {
			var row = []string{fmt.Sprintf("%v", s.ContestTime)}

			problem, hasProblem := problemSet(problems).byId(s.ProblemId)
			if hasProblem {
				row = append(row, fmt.Sprintf("%s: %s", problem.Label, problem.Name))
			} else {
				row = append(row, "unknown")
			}

			language, hasLanguage := languageSet(languages).byId(s.LanguageId)
			if hasLanguage {
				row = append(row, language.Name)
			} else {
				row = append(row, "unknown")
			}

			// Get judgement
			sjudgements, hasJudgements := judgementSet(judgements).bySubmissionId(s.Id)
			if hasJudgements {
				for _, j := range sjudgements {
					if j.JudgementTypeId != "" {
						judgementType, hasjudgementType := judgementTypeSet(judgementTypes).byId(j.JudgementTypeId)
						if hasjudgementType {
							row = append(row, fmt.Sprintf("%v", j.EndContestTime), fmt.Sprintf("%s (%s)", judgementType.Id, judgementType.Name))
						} else {
							row = append(row, fmt.Sprintf("%v", j.EndContestTime), "Unknown judgement")
						}
					} else {
						row = append(row, fmt.Sprintf("%v", j.StartContestTime), "In progress...")
					}
				} // todo handle multiple
			} else {
				row = append(row, "", "Queued")
			}
			table.appendRow(row)
		}
	}
	table.print()

	return nil
}
