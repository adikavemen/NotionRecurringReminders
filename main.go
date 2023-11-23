package main

import (
	"fmt"
	"strings"
	"os"
	"io"
	"net/http"
	"encoding/json"
	"time"
)

type Response struct {
	has_more bool `json:"has_more"`
	next_cursor string `json:"next_cursor"`
}

func errorHandle(err error){
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
		return
	}
}

func makeRequest(url string, startPoint string, notok string) []byte {
	var payload io.Reader
	if startPoint != "" {
		payload = strings.NewReader("{\"start_cursor\":\"" + startPoint + "\"}")
	}
	
	req, err := http.NewRequest("POST", url, payload)
	errorHandle(err)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer " + notok)

	res, err := http.DefaultClient.Do(req)
	errorHandle(err)

	if res.StatusCode != http.StatusOK{
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		errorHandle(err)
		return body
	}
	return nil
}

func main() {

	notionToken, notionExists := os.LookupEnv("NOTION_TOKEN")
	database, databaseExists := os.LookupEnv("DATABASE_ADDR")

	if notionExists && databaseExists {
		url := "https://api.notion.com/v1/databases/" + database + "/query"

		body := makeRequest(url, "", notionToken)
		
		fmt.Println(string(body))

		var data Response
		err := json.Unmarshal(body, &data)
		errorHandle(err)

		for data.has_more {
			time.Sleep(335 * time.Millisecond)
			if(data.has_more){
				body = makeRequest(url, data.next_cursor, notionToken)
				
				var data Response
				err = json.Unmarshal(body, &data)
				errorHandle(err)

				fmt.Println(string(body))
				
			}
		}
		

		
	} else if !notionExists {
		fmt.Println("Failed to find token")
		os.Exit(2)
	} else {
		fmt.Println("Failed to find database")
		os.Exit(3)
	}


	

}