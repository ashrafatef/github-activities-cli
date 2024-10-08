/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"githubActivitiesCli/database"
	"githubActivitiesCli/ui"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/go-git/go-git/v5"
	"github.com/guumaster/logsymbols"
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

var createPullRequestCmd = &cobra.Command{
	Use:     "create",
	Short:   "create pull request from current branch ",
	Run:     createPullRequest,
	Aliases: []string{"c"},
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
func runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

func getCurrentBranchName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	r, err := git.PlainOpen(dir)
	h, err := r.Head()
	// h, err := r.Log()
	currentBranch := h.Name().Short()
	// commit, err := r.CommitObject(h.Hash())
	// fmt.Println(commit)
	if err != nil {
		log.Fatal(err)
	}
	return currentBranch
}

func createPullRequest(cmd *cobra.Command, args []string) {

	currentBranch := getCurrentBranchName()

	err := runGitCommand("push", "-u", "origin", currentBranch)
	if err != nil {
		fmt.Println(logsymbols.Error, "Push branch failed")
		panic(err)
	}
	fmt.Println(logsymbols.Success, "Branch Pushed")

	token, err := database.GetToken()
	if err != nil {
		panic(err)
	}

	if len(token) == 0 {
		tokenPrompt := promptui.Prompt{
			Label: "Token",
			Mask:  '*',
		}
		promptToken, _ := tokenPrompt.Run()
		database.AddToken(promptToken)
		token = promptToken
	}

	titlePrompt := promptui.Prompt{
		Label: "title",
	}

	title, _ := titlePrompt.Run()
	prType := getPrType(cmd)
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	repo := path[strings.LastIndex(path, "/")+1:]

	var prTitle string

	if prType == "TECH" {
		prTitle = "feat: " + "TECH" + title + " " + title
	} else {
		prTitle = prType + ": " + currentBranch + " " + title
	}

	marshalled, err := json.Marshal(map[string]interface{}{
		"title": prTitle,
		"head":  currentBranch,
		"base":  "master",
		"body":  "",
	})

	fmt.Println(string(marshalled))

	url := "https://api.github.com/repos/ashrafatef/" + repo + "/pulls"
	authorization := "Bearer " + token
	client := http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalled))

	request.Header.Set("Authorization", authorization)
	request.Header.Set("Content-Type", "application/json")

	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	_, error := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusCreated {
		panic(fmt.Errorf("status code error: %d %s %s", res.StatusCode, res.Status, res))
	}

	if error != nil {
		panic(error)
	}
	fmt.Println(logsymbols.Success, "PR Created")

}

func getPrType(cmd *cobra.Command) string {
	ft, _ := cmd.Flags().GetBool("ft")
	fi, _ := cmd.Flags().GetBool("fi")
	t, _ := cmd.Flags().GetBool("t")

	if ft {
		return "feat"
	}
	if fi {
		return "fix"
	}
	if t {
		return "TECH"
	}
	panic("Please Provide pr type -ft for feat and -fi for fix and -t for TECH")
}

func init() {
	rootCmd.AddCommand(pullRequestsCmd)
	pullRequestsCmd.AddCommand(createPullRequestCmd)
	createPullRequestCmd.Flags().Bool("ft", false, "create feat PR")
	createPullRequestCmd.Flags().Bool("fi", false, "create fix PR")
	createPullRequestCmd.Flags().Bool("t", false, "create TECH PR")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullRequestsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullRequestsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
