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

	fmt.Printf("Submissions (%d):\n", count)
	for _, s := range submissions {
		if s.TeamId == teamId {
			fmt.Printf("  Submission to ")

			problem, hasProblem := problemSet(problems).byId(s.ProblemId)
			if hasProblem {
				fmt.Printf("problem %s: %s", problem.Label, problem.Name)
			} else {
				fmt.Printf("unknown problem")
			}

			language, hasLanguage := languageSet(languages).byId(s.LanguageId)
			if hasLanguage {
				fmt.Printf(" in %s", language.Name)
			} else {
				fmt.Printf(" in unknown language")
			}
			fmt.Printf(" at %s\n", s.ContestTime)

			// Get judgement
			sjudgements, hasJudgements := judgementSet(judgements).bySubmissionId(s.Id)
			if hasJudgements {
				for _, j := range sjudgements {
					if j.JudgementTypeId != "" {
						judgementType, hasjudgementType := judgementTypeSet(judgementTypes).byId(j.JudgementTypeId)
						if hasjudgementType {
							fmt.Printf("     Judged at %s: %s (%s)\n", j.EndContestTime, judgementType.Id, judgementType.Name)
						} else {
							fmt.Printf("     Unknown judgement at %s\n", j.EndContestTime)
						}
					} else {
						fmt.Printf("     Judgement in progress\n")
					}
				}
			} else {
				fmt.Println("     Not judged")
			}
		}
	}

	return nil
}
