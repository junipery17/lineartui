package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func testAddIssue(t *testing.T) {
	url := "https://api.linear.app/graphql"
	key := os.Getenv("linearAPI_Key")
	client := &http.Client{}

	resp, err := AddIssue(client, url, key, "ec0b720c-d68c-4907-9708-8a3e52b810cc", "test")

	if (err != nil) || (resp == "") {
		t.Errorf(`Something is wrong with connecting with the API me thinks: %q`, err)
	}
	var returned Response
	json.Unmarshal([]byte(resp), &returned)
	if returned.Data.IssueCreate.Success == false {
		t.Error("problem :( the issue was not created")
	}
}

func testAddIssueWithDes(t *testing.T) {
	url := "https://api.linear.app/graphql"
	key := os.Getenv("linearAPI_Key")
	client := &http.Client{}

	resp, err := AddIssue(client, url, key, "ec0b720c-d68c-4907-9708-8a3e52b810cc", "test", "This is the description yup mhm")

	if (err != nil) || (resp == "") {
		t.Errorf(`Something is wrong with connecting with the API me thinks: %q`, err)
	}
	var returned Response
	json.Unmarshal([]byte(resp), &returned)
	if returned.Data.IssueCreate.Success == false {
		t.Error("problem :( the issue was not created")
	}
}

func testAddIssueWithNoTitle(t *testing.T) {
	url := "https://api.linear.app/graphql"
	key := os.Getenv("linearAPI_Key")
	client := &http.Client{}

	resp, err := AddIssue(client, url, key, "ec0b720c-d68c-4907-9708-8a3e52b810cc", "")

	if (err != nil) || (resp == "") {
		t.Errorf(`Something is wrong with connecting with the API me thinks: %q`, err)
	}
	var returned Response
	json.Unmarshal([]byte(resp), &returned)
	if returned.Data.IssueCreate.Success == false {
		t.Error("problem :( the issue was not created")
	}
}

func testDeleteIssue(t *testing.T) {
	url := "https://api.linear.app/graphql"
	key := os.Getenv("linearAPI_Key")
	client := &http.Client{}

	resp, err := AddIssue(client, url, key, "ec0b720c-d68c-4907-9708-8a3e52b810cc", "test")

	if (err != nil) || (resp == "") {
		t.Errorf(`Something is wrong with connecting with the API me thinks: %q`, err)
	}
	var returned Response
	json.Unmarshal([]byte(resp), &returned)
	if returned.Data.IssueCreate.Success == false {
		t.Error("problem :( the issue was not created")
	}

	id := returned.Data.IssueCreate.Issue.Id
	DeleteIssue(client, url, key, id)
}
