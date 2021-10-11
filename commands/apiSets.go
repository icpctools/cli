package commands

import (
	"strings"

	interactor "github.com/icpctools/api-interactor"
)

type (
	judgementTypeSet []interactor.JudgementType
	problemSet       []interactor.Problem
	languageSet      []interactor.Language
	judgementSet     []interactor.Judgement
	teamSet          []interactor.Team
)

func (j judgementTypeSet) byId(id string) (interactor.JudgementType, bool) {
	for _, jt := range j {
		if jt.Id == id {
			return jt, true
		}
	}

	return interactor.JudgementType{}, false
}

func (t teamSet) byId(id string) (interactor.Team, bool) {
	for _, team := range t {
		if team.Id == id {
			return team, true
		}
	}

	return interactor.Team{}, false
}

func (j judgementSet) bySubmissionId(id string) (judgementSet, bool) {
	var jud []interactor.Judgement
	for _, ju := range j {
		if ju.SubmissionId == id {
			jud = append(jud, ju)
		}
	}

	return jud, len(jud) > 0
}

func (p problemSet) byId(id string) (interactor.Problem, bool) {
	for _, problem := range p {
		if strings.EqualFold(problem.Id, id) || strings.EqualFold(problem.Label, id) || strings.EqualFold(problem.Name, id) {
			return problem, true
		}
	}

	return interactor.Problem{}, false
}

func (l languageSet) byId(id string) (interactor.Language, bool) {
	for _, language := range l {
		if language.Id == id {
			return language, true
		}
	}

	return interactor.Language{}, false
}

func (l languageSet) byExtension(extension string) (interactor.Language, bool) {
	for _, language := range l {
		for _, languageExtension := range language.Extensions {
			if strings.ToLower(languageExtension) == extension {
				return language, true
			}
		}
	}

	return interactor.Language{}, false
}
