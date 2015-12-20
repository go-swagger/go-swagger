// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

package main

import (
	"flag"
	"github.com/golang/gddo/gosrc"
	"log"
	"time"
)

var backgroundTasks = []*struct {
	name     string
	fn       func() error
	interval *time.Duration
	next     time.Time
}{
	{
		name:     "GitHub updates",
		fn:       readGitHubUpdates,
		interval: flag.Duration("github_interval", 0, "Github updates crawler sleeps for this duration between fetches. Zero disables the crawler."),
	},
	{
		name:     "Crawl",
		fn:       doCrawl,
		interval: flag.Duration("crawl_interval", 0, "Package updater sleeps for this duration between package updates. Zero disables updates."),
	},
}

func runBackgroundTasks() {
	defer log.Println("ERROR: Background exiting!")

	sleep := time.Minute
	for _, task := range backgroundTasks {
		if *task.interval > 0 && sleep > *task.interval {
			sleep = *task.interval
		}
	}

	for {
		for _, task := range backgroundTasks {
			start := time.Now()
			if *task.interval > 0 && start.After(task.next) {
				if err := task.fn(); err != nil {
					log.Printf("Task %s: %v", task.name, err)
				}
				task.next = time.Now().Add(*task.interval)
			}
		}
		time.Sleep(sleep)
	}
}

func doCrawl() error {
	// Look for new package to crawl.
	importPath, hasSubdirs, err := db.PopNewCrawl()
	if err != nil {
		log.Printf("db.PopNewCrawl() returned error %v", err)
		return nil
	}
	if importPath != "" {
		if pdoc, err := crawlDoc("new", importPath, nil, hasSubdirs, time.Time{}); pdoc == nil && err == nil {
			if err := db.AddBadCrawl(importPath); err != nil {
				log.Printf("ERROR db.AddBadCrawl(%q): %v", importPath, err)
			}
		}
		return nil
	}

	// Crawl existing doc.
	pdoc, pkgs, nextCrawl, err := db.Get("-")
	if err != nil {
		log.Printf("db.Get(\"-\") returned error %v", err)
		return nil
	}
	if pdoc == nil || nextCrawl.After(time.Now()) {
		return nil
	}
	if _, err = crawlDoc("crawl", pdoc.ImportPath, pdoc, len(pkgs) > 0, nextCrawl); err != nil {
		// Touch package so that crawl advances to next package.
		if err := db.SetNextCrawlEtag(pdoc.ProjectRoot, pdoc.Etag, time.Now().Add(*maxAge/3)); err != nil {
			log.Printf("ERROR db.TouchLastCrawl(%q): %v", pdoc.ImportPath, err)
		}
	}
	return nil
}

func readGitHubUpdates() error {
	const key = "gitHubUpdates"
	var last string
	if err := db.GetGob(key, &last); err != nil {
		return err
	}
	last, names, err := gosrc.GetGitHubUpdates(httpClient, last)
	if err != nil {
		return err
	}

	for _, name := range names {
		log.Printf("bump crawl github.com/%s", name)
		if err := db.BumpCrawl("github.com/" + name); err != nil {
			log.Println("ERROR force crawl:", err)
		}
	}

	if err := db.PutGob(key, last); err != nil {
		return err
	}
	return nil
}
