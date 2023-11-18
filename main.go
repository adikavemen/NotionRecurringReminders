package main

import (
	"fmt"
	"strings"
	"os"
	"io"
	"net/http"
)

func main() {

	notionToken, exists := os.LookupEnv("NOTION_TOKEN")

	if exists {
		url := "https://api.notion.com/v1/databases/04e8df4e-60bd-4f2d-83ad-0e4edcb0941e/query"

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