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

type GithubEventsResponse struct {
	ID      string `json:"id"`
	Payload struct {
		Issue struct {
			CreatedAt string `json:"created_at"`
		} `json:"issue"`
		Comment struct {
			CreatedAt string `json:"created_at"`
		} `json:"comment"`
		PullRequest struct {
			CreatedAt string `json:"created_at"`
			Title     string `json:"title"`
			Head      struct {
				Ref  string `json:"ref"`
				Repo struct {
					Name string `json:"name"`
				} `json:"repo"`
			} `json:"head"`
		} `json:"pull_request"`
		Review struct {
			CreatedAt string `json:"created_at"`
		} `json:"review"`
		Commits []struct {
			CreatedAt string `json:"created_at"`
			Message   string `json:"message"`
		} `json:"commits"`
	}
}

// pullRequestsCmd represents the pullRequests command
var pullRequestsCmd = &cobra.Command{
	Use:   "pullRequests",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run:     pullRequests,
	Aliases: []string{"pr"},
}

func pullRequests(cmd *cobra.Command, args []string) {
	usernamePrompt := promptui.Prompt{
		Label: "Username",
	}

	tokenPrompt := promptui.Prompt{
		Label: "Token",
		Mask:  '*',
	}

	username, _ := usernamePrompt.Run()
	token, _ := tokenPrompt.Run()
	url := "https://api.github.com/users/" + username + "/events"
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

	var githubResponses []GithubEventsResponse
	_ = json.Unmarshal(body, &githubResponses)
	var rows []table.Row

	for _, githubResponse := range githubResponses {
		if githubResponse.Payload.PullRequest.CreatedAt != "" {
			rows = append(rows, table.Row{
				githubResponse.Payload.PullRequest.CreatedAt,
				githubResponse.Payload.PullRequest.Title,
				githubResponse.Payload.PullRequest.Head.Ref,
				githubResponse.Payload.PullRequest.Head.Repo.Name,
			})
		}
	}

	ui.RunProgress()
	columns := []table.Column{
		{Title: "CreatedAt", Width: 20},
		{Title: "Title", Width: 70},
		{Title: "Branch Name", Width: 20},
		{Title: "Repo Name", Width: 20},
	}

	ui.CreateTable(rows, columns)

}

func init() {
	rootCmd.AddCommand(pullRequestsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullRequestsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullRequestsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
