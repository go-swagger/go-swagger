//go:build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	allVersions bool
	version     = "latest"
)

func init() {
	flag.BoolVar(&allVersions, "all", allVersions, "when specified it will download stats for all versions")
	flag.StringVar(&version, "version", version, "the version to download stats for")
}

func main() {
	flag.Parse()

	req, err := http.NewRequest("GET", "https://api.github.com/repos/go-swagger/go-swagger/releases", nil)
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
		b, err := io.ReadAll(resp.Body)
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
		if !allVersions {
			if (version == "latest" || version == "") && i > 0 {
				break
			}
		}
		if allVersions && i > 0 {
			fmt.Println()
		}

		if !allVersions && (version != "latest" && version != "") && result.TagName != version {
			continue
		}

		fmt.Println("Stats for release:", result.TagName)
		for _, asset := range result.Assets {
			fmt.Printf("%25s: %d\n", asset.Name, asset.Downloads)
		}

		if !allVersions && (version != "latest" && version != "") && result.TagName != version {
			break
		}
	}
}
