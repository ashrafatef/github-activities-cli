/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"githubActivitiesCli/ui"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/table"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Workflow struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	DisplayTitle string `json:"display_title"`
}

type GitHubWorkflows struct {
	TotalCount   int32      `json:"total_count"`
	WorkflowRuns []Workflow `json:"workflow_runs"`
}

// workflowsCmd represents the workflows command
var workflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"wo"},
	Run: func(cmd *cobra.Command, args []string) {
		workflowName, _ := cmd.Flags().GetString("name")
		tokenPrompt := promptui.Prompt{
			Label: "Token",
			Mask:  '*',
		}

		repoPrompt := promptui.Prompt{
			Label: "Repo",
		}

		token, _ := tokenPrompt.Run()
		repo, _ := repoPrompt.Run()

		url := "https://api.github.com/repos/join-com/" + repo + "/actions/runs"
		authorization := "Bearer " + token
		client := http.Client{}
		request, _ := http.NewRequest(http.MethodGet, url, nil)

		request.Header.Set("Authorization", authorization)
		request.Header.Set("Content-Type", "application/json")

		res, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			panic(fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status))
		}

		body, error := io.ReadAll(res.Body)
		if error != nil {
			panic(error)
		}

		var githubResponses GitHubWorkflows
		_ = json.Unmarshal(body, &githubResponses)
		var rows []table.Row

		for _, workflows := range githubResponses.WorkflowRuns {
			
			if !(workflowName != "" && workflows.Name == workflowName) {
				continue
			}
			rows = append(rows, table.Row{
				string(workflows.ID),
				workflows.Name,
				workflows.Status,
				workflows.DisplayTitle,
			})
		}

		ui.RunProgress()
		columns := []table.Column{
			{Title: "ID", Width: 20},
			{Title: "Name", Width: 20},
			{Title: "Status", Width: 20},
			{Title: "Display Title", Width: 70},
		}

		ui.CreateTable(rows, columns)
	},
}

func init() {
	rootCmd.AddCommand(workflowsCmd)
	workflowsCmd.PersistentFlags().String("name", "n", "workflow name")
}
