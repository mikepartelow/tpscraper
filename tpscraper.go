package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const URL_TEMPLATE = "https://listen.tidal.com/v1/playlists/%s/items?offset=%d&limit=%d&countryCode=US&locale=en_US&deviceType=BROWSER"
const LIMIT = 100

type Artist struct {
	Name string `json:"name"`
}

type Details struct {
	Title  string `json:"title"`
	Artist Artist `json:"artist"`
}

type Item struct {
	Type    string  `json:"type"`
	Details Details `json:"item"`
}

type Items struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Total  int    `json:"totalNumberOfItems"`
	Items  []Item `json:"items"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseJson(theJson []byte) (items Items) {
	err := json.Unmarshal(theJson, &items)
	check(err)

	return
}

func fetchJson(playlist_id string, auth_token string) []byte {
	url := fmt.Sprintf(URL_TEMPLATE, playlist_id, 0, LIMIT)

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	check(err)

	request.Header.Set("User-Agent", "Not Firefox")
	request.Header.Set("Authorization", auth_token)

	response, err := client.Do(request)
	check(err)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	check(err)

	return body
}

func readSecrets() (string, string) {
	playlist_id, err := os.ReadFile("playlist_id")
	check(err)

	auth_token, err := os.ReadFile("auth_token")
	check(err)

	return string(playlist_id), string(auth_token)
}

func main() {
	playlist_id, auth_token := readSecrets()

	theJson := fetchJson(playlist_id, auth_token)

	items := parseJson(theJson)

	fmt.Printf("Total: %d, Batch: %d\n", items.Total, len(items.Items))

	for _, item := range items.Items {
		fmt.Printf("%s: %s\n", item.Details.Artist.Name, item.Details.Title)
	}
}
