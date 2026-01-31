package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	TeamNodes   Nodes       `json:"teams"`
	Team        TeamData    `json:"team"`
	IssueCreate IssueCreate `json:"issueCreate"`
}

type IssueCreate struct {
	Success bool      `json:"success"`
	Issue   IssueData `json:"issue"`
}

type Nodes struct {
	EachTeam []TeamData `json:"nodes"`
}

type Issues struct {
	EachIssue []IssueData `json:"nodes"`
}

type TeamData struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Issues Issues `json:"issues"`
}

type IssueData struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Assignee    struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"assignee"`
}

func postQueryRequest(client *http.Client, url string, key string, query string) string {
	jsonData := map[string]string{
		"query": query,
	}

	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("uh ohh")
	}
	request.Header.Add("Authorization", key)
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return string(body)
}

func DisplayTeams(client *http.Client, url string, key string) {
	query := `
		query Teams {
			teams {
				nodes {
					id
					name
				}
			}
		}
	`
	sb := postQueryRequest(client, url, key, query)
	var teams_data Response
	json.Unmarshal([]byte(sb), &teams_data)
	length := len(teams_data.Data.TeamNodes.EachTeam)
	for i := range length {
		fmt.Printf("Team: %#v\n", teams_data.Data.TeamNodes.EachTeam[i])
	}
}

func DisplayIssues(client *http.Client, url string, key string, teamId string) {
	//theres like a lot of things that can be added to display so idk which ones are like important
	query := fmt.Sprintf(`query Team {
			team(id: %q) {
				name
				issues {
					nodes {
						id
						title
						description
						assignee {
						id
						name
						}
					}
				}
			}
			}`, teamId)

	sb := postQueryRequest(client, url, key, query)
	var issues_data Response
	json.Unmarshal([]byte(sb), &issues_data)
	length := len(issues_data.Data.Team.Issues.EachIssue)
	for i := range length {
		fmt.Printf("%#v\n", issues_data.Data.Team.Issues.EachIssue[i])
	}
}

func AddIssue(client *http.Client, url string, key string, teamId string, title string, description ...string) (string, error) {
	if title == "" {
		return "", errors.New("no title")
	}
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	mutation := fmt.Sprintf(`
		mutation IssueCreate{
			issueCreate(
				input: {
					title: %q
					description: %q
					teamId: %q
				}
			)
			{
				success
				issue{
					id
					title
				}
			}
		}
	`, title, desc, teamId)
	sb := postQueryRequest(client, url, key, mutation)
	fmt.Println(sb)
	return sb, nil
}

func DeleteIssue(client *http.Client, url string, key string, issueID string) (string, error) {
	mutation := fmt.Sprintf(`
		mutation IssueDelete{
			issueDelete(
				id: %q
			) {
			success
		}
		}
	`, issueID)
	sb := postQueryRequest(client, url, key, mutation)
	fmt.Println(string(sb))
	return sb, nil
}

func main() {
	url := "https://api.linear.app/graphql"
	key := os.Getenv("linearAPI_Key")
	client := &http.Client{}
	// display teams
	DisplayTeams(client, url, key)

	// display issues of team
	DisplayIssues(client, url, key, "ec0b720c-d68c-4907-9708-8a3e52b810cc")

	//Create an Issue without description
	// AddIssue(client, url, key, "ec0b720c-d68c-4907-9708-8a3e52b810cc", "testing")

	//Create an Issue with description
	// AddIssue(client, url, key, "ec0b720c-d68c-4907-9708-8a3e52b810cc", "testing", "with a description")

	//Delete an Issue
	// DeleteIssue(client, url, key, "")

}

// "ec0b720c-d68c-4907-9708-8a3e52b810cc" <- id of yuzubox team
