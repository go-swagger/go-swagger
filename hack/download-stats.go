package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/go-swagger/go-swagger/releases", nil)
	// req, err := http.NewRequest("GET", "https://api.github.com/repos/go-swagger/go-swagger/releases/latest", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err) //nolint: gocritic
		}
		log.Fatalf("%s: %s", resp.Status, b) //nolint: gocritic
	}

	var results []struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name      string `json:"name"`
			Downloads int64  `json:"download_count"`
		} `json:"assets"`
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&results); err != nil {
		log.Fatalln(err)
	}

	for i, result := range results {
		if i > 0 {
			fmt.Println()
		}

		fmt.Println("Stats for release:", result.TagName)
		for _, asset := range result.Assets {
			fmt.Printf("%25s: %d\n", asset.Name, asset.Downloads)
		}
	}
}
