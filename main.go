package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var url string

func errorHandle(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
		return
	}
}

// Retrives the set of tasks from the notion database where the tasks are created.
// startPoint - The address of the current set to retrieve
// notok - notion token
func getTaskSet(startPoint string, notok string) []byte {
	var payload io.Reader
	if startPoint != "" {
		payload = strings.NewReader("{\"start_cursor\":\"" + startPoint + "\"}")
	} else {
		payload = strings.NewReader("{}")
	}

	req, err := http.NewRequest("POST", url, payload)
	errorHandle(err)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+notok)

	res, err := http.DefaultClient.Do(req)
	errorHandle(err)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil
	} else {
		body, err := io.ReadAll(res.Body)
		errorHandle(err)
		return body
	}
}

func getAllTasks(notionToken string, database string) {
	url = "https://api.notion.com/v1/databases/" + database + "/query"

	body := getTaskSet("", notionToken)

	// TODO: Do something with the tasks results here
	fmt.Println(string(body))

	var data ResultsQuery
	err := json.Unmarshal(body, &data)
	errorHandle(err)

	for data.has_more {
		time.Sleep(335 * time.Millisecond)
		if data.has_more {
			body = getTaskSet(data.next_cursor, notionToken)

			var data ResultsQuery
			err = json.Unmarshal(body, &data)
			errorHandle(err)

			// TODO: Do something with the tasks results here
			fmt.Println(string(body))
		}
	}
}

func main() {
	notionToken, notionExists := os.LookupEnv("NOTION_TOKEN")
	database, databaseExists := os.LookupEnv("DATABASE_ADDR")

	if notionExists && databaseExists {
		getAllTasks(notionToken, database)
	} else {
		content, err := os.ReadFile(".env")
		errorHandle(err)
		lines := strings.Split(string(content), "\n")
		notionTok := lines[0]
		dataAddr := lines[1]
		getAllTasks(notionTok, dataAddr)
	}

	// else if !notionExists {
	// 	fmt.Println("Failed to find token")
	// 	os.Exit(2)
	// } else {
	// 	fmt.Println("Failed to find database")
	// 	os.Exit(3)
	// }

}
