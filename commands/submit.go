package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	interactor "github.com/tuupke/api-interactor"
)

var submitCommand = &cobra.Command{
	Use:     "submit [file1] <file2> <file3> ...",
	Short:   "Submit one or more files",
	Args:    cobra.MinimumNArgs(1),
	RunE:    submit,
	PreRunE: configHelper("baseurl"),
}

type (
	problemSet  []interactor.Problem
	languageSet []interactor.Language
)

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

func submit(cmd *cobra.Command, args []string) error {
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	contest, err := api.ContestById(viper.GetString("contest"))
	if err != nil {
		return fmt.Errorf("could not get contest; %w", err)
	}

	// Get the problems and languages
	problems, err := api.Problems()
	if err != nil {
		return fmt.Errorf("could not get problems; %w", err)
	}

	languages, err := api.Languages()
	if err != nil {
		return fmt.Errorf("could not get languages; %w", err)
	}

	// Try to load all the files from the arguments
	var files interactor.LocalFileReference
	for _, filename := range args {
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("could not open file %s; %w", filename, err)
		}

		err = files.FromFile(file)
		if err != nil {
			return fmt.Errorf("could not add file %s; %w", filename, err)
		}
	}

	// If the problem or language is not set, use the first file to determine them
	if problemId == "" || languageId == "" {
		// Assume first part of the basename can be used to detect problem and the extension can be used to detect language
		firstFileParts := strings.Split(filepath.Base(args[0]), ".")
		var extension string
		if len(firstFileParts) > 1 {
			extension = strings.ToLower(firstFileParts[len(firstFileParts)-1])
		}

		if problemId == "" {
			problemId = strings.ToLower(firstFileParts[0])
		}

		if languageId == "" {
			if language, found := languageSet(languages).byExtension(extension); found {
				languageId = language.Id
			}
		}
	}

	problem, hasProblem := problemSet(problems).byId(problemId)
	language, hasLanguage := languageSet(languages).byId(languageId)

	if !hasProblem {
		return fmt.Errorf("no known problem specified or detected")
	}

	if !hasLanguage {
		return fmt.Errorf("no known language specified or detected")
	}

	// Try to auto detect entry point based on hardcoded language logic
	if entryPoint == "" && language.EntryPointRequired {
		switch language.Id {
		case "java":
			// Java: use base name of first file
			parts := strings.Split(filepath.Base(args[0]), ".")
			entryPoint = parts[0]
		case "python", "python2", "python3":
			// Python: use first file
			entryPoint = filepath.Base(args[0])
		case "kotlin":
			parts := strings.Split(filepath.Base(args[0]), ".")
			entryPoint = kotlinBaseEntryPoint(parts[0]) + "Kt"
		}
	}

	if entryPoint == "" && language.EntryPointRequired {
		return fmt.Errorf("entry point required but not specified nor detected")
	}

	if !force {
		fmt.Println("About to submit:")
		if len(args) == 1 {
			fmt.Printf("  filename:    %s\n", args[0])
		} else {
			fmt.Print("  filenames:  ")
			for _, filename := range args {
				fmt.Printf(" %s\n", filename)
			}
		}

		fmt.Printf("  contest:     %s\n", contest.Name)
		fmt.Printf("  problem:     %s\n", problem.Label)
		fmt.Printf("  language:    %s\n", language.Name)
		if entryPoint != "" {
			fmt.Printf("  entry point: %s\n", entryPoint)
		}

		if !prompter.YN("Do you want to submit?", true) {
			return errors.New("submission aborted by user")
		}
	}

	submissionId, err := api.PostSubmission(problem.Id, language.Id, entryPoint, files)
	if err != nil {
		return fmt.Errorf("could not submit: %w", err)
	}

	fmt.Println("Submitted. ID:", submissionId)
	return nil
}

func kotlinBaseEntryPoint(base string) string {
	if base == "" {
		return "_"
	}

	out := []rune(base)
	for i, r := range out {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			out[i] = '_'
		}
	}

	if unicode.IsLetter(out[0]) {
		out[0] = unicode.ToUpper(out[0])
		return string(out)
	}

	return "_" + string(out)
}
