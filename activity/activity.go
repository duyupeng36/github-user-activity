package activity

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var urlTemp = "https://api.github.com/users/%s/events"

type Repo struct {
	Name string `json:name`
}
type GithubActivity struct {
	Type      string `json:"type"`
	Repo      Repo   `json:"repo"`
	CreatedAt string `json:"created_at"`
	Payload   struct {
		Action  string `json:"action"`
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
	} `json:"payload"`
}

func request(username string, perPage, page int) (*http.Request, error) {

	url := fmt.Sprintf(urlTemp, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.URL.Query().Set("per_page", strconv.Itoa(perPage))
	req.URL.Query().Set("page", strconv.Itoa(page))
	return req, nil
}

func getActivities(username string, perPage, page int) ([]GithubActivity, error) {
	req, err := request(username, perPage, page)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found. please check the username")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching data: %d", resp.StatusCode)
	}

	// 解析为json
	var activities = make([]GithubActivity, 0)
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, err
	}
	return activities, nil
}

func ListAllActivities(username string, perPage, page int) {
	log.Printf("Starting to obtain %d activities for %d pages of %s in total\n", perPage, page, username)
	activities, err := getActivities(username, perPage, page)

	if err != nil {
		log.Printf("Error of obtaining activities: %v\n", err)
	}

	log.Printf("Found %d activities\n", len(activities))

	for _, activity := range activities {
		var action string
		switch activity.Type {
		case "PushEvent":
			commitCount := len(activity.Payload.Commits)
			action = fmt.Sprintf("Pushed %d commit(s) to %s", commitCount, activity.Repo.Name)
		case "IssueCommentEvent":
			action = fmt.Sprintf("%s an issue in %s", activity.Payload.Action, activity.Repo.Name)
		case "WatchEvent":
			action = fmt.Sprintf("Starred %s", activity.Repo.Name)
		case "ForkEvent":
			action = fmt.Sprintf("Forked %s", activity.Repo.Name)
		case "CreateEvent":
			action = fmt.Sprintf("Created %s in %s", activity.Payload.RefType, activity.Repo.Name)
		default:
			action = fmt.Sprintf("%s in %s", activity.Type, activity.Repo.Name)
		}
		fmt.Println("-", action)
	}
}
