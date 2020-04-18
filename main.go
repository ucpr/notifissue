package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Issue struct {
	Title     string `json:"title"`
	UpdatedAt string `json:"updated_at"`
}

type PullRequest struct {
	Title     string `json:"title"`
	UpdatedAt string `json:"updated_at"`
}

type Event struct {
	EventType string `json:"type"`
	Payload   struct {
		Action      string      `json:"action"`
		Issue       Issue       `json:"issue"`
		PullRequest PullRequest `json:"pull_request"`
	} `json:"payload"`
}

func main() {
	var username string = parseArgs()
	bytes, err := fetchEvents(username)
	if err != nil {
		log.Fatal(err)
	}

	var events []Event
	if err := json.Unmarshal(bytes, &events); err != nil {
		log.Fatal(err)
	}

	fmt.Print("GitHub Recent activities\n\n")
	printEvents(events)
}

func parseArgs() string {
	var username = flag.String("u", "", "username")
	flag.Parse()
	return *username
}

func printEvents(events []Event) {
	for _, event := range events {
		switch event.EventType {
		case "IssuesEvent":
			printIssue(event, "[Issue]")
		case "PullRequestEvent":
			printPullRequest(event, "[PullRequest]")
		}
	}
}

func printIssue(event Event, eventType string) {
	fmt.Printf(
		"%s %s\n%s\n\n",
		event.Payload.Action,
		eventType,
		event.Payload.Issue.Title,
	)
}

func printPullRequest(event Event, eventType string) {
	fmt.Printf(
		"%s %s\n%s\n\n",
		event.Payload.Action,
		eventType,
		event.Payload.PullRequest.Title,
	)
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func fetchEvents(username string) ([]byte, error) {
	var url string = "https://api.github.com/users/" + username + "/events/public"
	return fetch(url)
}
