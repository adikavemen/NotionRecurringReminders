package main

import (
	"fmt"
	"strings"
	"os"
	"io"
	"net/http"
)

func main() {

	notionToken, notionExists := os.LookupEnv("NOTION_TOKEN")
	database, databaseExists := os.LookupEnv("DATABASE_ADDR")

	if notionExists && databaseExists {
		url := "https://api.notion.com/v1/databases/" + database + "/query"

		payload := strings.NewReader("{\"filter\": {\"and\": [{\"property\": \"Status\",\"status\": {\"equals\": \"Not Started\"}},{\"property\": \"Tags\",\"multi_select\": {\"contains\": \"Daily\"}}]}}")

		req, _ := http.NewRequest("POST", url, payload)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("User-Agent", "Insomnia/2023.5.6")
		req.Header.Add("Notion-Version", "2022-06-28")
		req.Header.Add("Authorization", "Bearer " + notionToken)

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		fmt.Println(string(body))	
	} else {
		fmt.Println("Failed to find token")
	}


	

}