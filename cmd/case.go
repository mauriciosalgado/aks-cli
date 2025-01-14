/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// Define a hardcoded list of tags
var availableTags = []string{
	"#bug", "#enhancement", "#urgent", "#help", "#feature-request", "#discussion", "#wontfix",
}

func init() {
	// Here you will define your flags and configuration settings.
	searchCmd.AddCommand(caseCmd)
	caseCmd.Flags().BoolP("tags", "t", false, "Search cases by tag(s)")
}

// casesCmd represents the cases command
var caseCmd = &cobra.Command{
	Use:   "case",
	Short: "Searches for cases",
	Long:  `Enables filtering cases with or without criteria`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caseName := args[0]

		// Get the value of the "tags" flag (whether it was set to true)
		isTagsFlagSet, _ := cmd.Flags().GetBool("tags")

		// If the "tags" flag is set, launch fzf to select tags
		if isTagsFlagSet {
			selectedTags := selectTagsFZF(availableTags)
			fmt.Println("Case Name:", caseName)
			fmt.Println("Selected Tags:", strings.Join(selectedTags, " "))
		} else {
			// If no tags flag is set, proceed with the case name and no tag selection
			fmt.Println("Case Name:", caseName)
			fmt.Println("No tags selected.")
		}
	},
}

func selectTagsFZF(tags []string) []string {
	tagList := strings.Join(tags, "\n")

	// Run fzf with the tag list, capture the selected output
	cmd := exec.Command("fzf", "--multi", "--preview", "echo {}", "--height", "40%", "--border")
	cmd.Stdin = strings.NewReader(tagList)

	// Get the selected tags as output
	selectedTags, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running fzf:", err)
		return nil
	}

	// Clean the output by trimming spaces and split by newline to get an array of tags
	selectedTagsString := strings.TrimSpace(string(selectedTags))
	if selectedTagsString == "" {
		return nil
	}

	// Return the selected tags as a string array
	return strings.Split(selectedTagsString, "\n")
}

func selectCasesFZF(files []string) ([]string, error) {
	cmd := exec.Command("fzf", "multi", "--preview", "cat {}", "--prompt", "Select Cases > ")
	cmd.Stdin = strings.NewReader(strings.Join(files, "\n"))

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	selected := strings.Split(strings.TrimSpace(string(output)), "\n")

	return selected, nil
}

func openInNeovim(files []string) error {
	args := append([]string{}, files...)
	cmd := exec.Command("nvim", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
