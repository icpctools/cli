package commands

import (
	"errors"
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
	interactor "github.com/tuupke/api-interactor"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func init() {
	cmd := &cobra.Command{
		Use:   "submit [file1] <file2> <file3> ...",
		Short: "Submit one or more files",
		Args:  cobra.MinimumNArgs(1),
		RunE:  submit,
	}

	cmd.Flags().StringVar(&problemId, "problem", "", "problem ID to submit for. Leave empty to auto detect from first file")
	cmd.Flags().StringVarP(&languageId, "language", "l", "", "language ID to submit for. Leave empty to auto detect from first file")
	cmd.Flags().StringVarP(&entryPoint, "entry-point", "e", "", "entry point to use. Leave empty if not needed or to auto detect")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "whether to force submission (i.e. not ask for confirmation")

	rootCommand.AddCommand(cmd)
}

func submit(cmd *cobra.Command, args []string) error {
	if baseUrl == "" {
		return errors.New("no base URL provided in flag or config")
	}
	if contestId == "" {
		return errors.New("no contest ID provided in flag or config")
	}

	api, err := interactor.ContestInteractor(baseUrl, username, password, contestId, insecure)
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	contest, err := api.ContestById(contestId)
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
	files := interactor.NewLocalFileReference()
	for _, filename := range args {
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("could not open file %s; %w", filename, err)
		}

		err = files.AddFromFile(file)
		if err != nil {
			return fmt.Errorf("could not add file %s; %w", filename, err)
		}
	}

	// If the problem or language is not set, use the first file to determine them
	if problemId == "" || languageId == "" {
		// Assume first part of the basename can be used to detect problem and the extension can be used to detect language
		firstFileParts := strings.Split(filepath.Base(args[0]), ".")
		extension := ""
		if len(firstFileParts) > 1 {
			extension = strings.ToLower(firstFileParts[len(firstFileParts)-1])
		}

		if problemId == "" {
			problemId = strings.ToLower(firstFileParts[0])
		}

		if languageId == "" {
		languageLoop:
			for _, language := range languages {
				for _, languageExtension := range language.Extensions {
					if strings.ToLower(languageExtension) == extension {
						languageId = language.Id
						break languageLoop
					}
				}
			}
		}
	}

	var problem *interactor.Problem
	var language *interactor.Language

	for _, p := range problems {
		if strings.ToLower(p.Id) == strings.ToLower(problemId) || strings.ToLower(p.Label) == strings.ToLower(problemId) {
			problem = &p
			break
		}
	}

	for _, l := range languages {
		if l.Id == languageId {
			language = &l
			break
		}
	}

	if problem == nil {
		return fmt.Errorf("no known problem specified or detected")
	}

	if language == nil {
		return fmt.Errorf("no known language specified or detected")
	}

	// Try to auto detect entry point based on hardcoded language logic
	if entryPoint == "" && language.EntryPointRequired {
		switch language.Id {
		case "java":
			// Java: use base name of first file
			parts := strings.Split(filepath.Base(args[0]), ".")
			entryPoint = parts[0]
		case "python2":
		case "python3":
		case "python":
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
				fmt.Printf(" %s", filename)
			}
			fmt.Println()
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Submitted. ID: %s\n", submissionId)
	return err
}

func kotlinBaseEntryPoint(base string) string {
	if base == "" {
		return "_"
	}

	isAlpha := func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
	}
	isNum := func(r rune) bool {
		return r >= '0' && r <= '9'
	}

	out := []rune(base)
	for i, r := range out {
		if !isAlpha(r) && !isNum(r) {
			out[i] = '_'
		}
	}

	if isAlpha(out[0]) {
		out[0] = unicode.ToUpper(out[0])
		return string(out)
	} else {
		return "_" + string(out)
	}
}
