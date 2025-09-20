package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
)

func strptr(s string) *string {
	return &s
}

func chk(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func headersString(headers http.Header) string {
	var lines []string
	for key, values := range headers {
		for _, value := range values {
			lines = append(lines, fmt.Sprintf("%s: %s", key, value))
		}
	}
	return strings.Join(lines, "\n")
}

func doRequest(task Task, httpClient *http.Client, userAgent string) *Result {
	r := &Result{Id: task.Id}

	req, err := http.NewRequest("GET", task.Url, nil)
	if err != nil {
		r.Err = strptr(err.Error())
		return r
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := httpClient.Do(req)
	if err != nil {
		r.Err = strptr(err.Error())
		return r
	}
	defer res.Body.Close()

	r.Status = &res.StatusCode
	r.Headers = strptr(headersString(res.Header))
	r.Content, err = io.ReadAll(res.Body)

	if err != nil {
		r.Err = strptr(err.Error())
	}
	return r
}

func extractUrls(inputsFile string) []string {
	file, err := os.Open(inputsFile)
	chk(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var urls []string

	urlRegex := regexp.MustCompile(`https?://[^\s]+`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := urlRegex.FindAllString(line, -1)
		urls = append(urls, matches...)
	}

	chk(scanner.Err())
	return urls
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func main() {
	/////////////////
	// parse args

	cfg := parseArgs()

	if cfg.DbFile == "" {
		fmt.Println("ERR: db file is required")
		flag.Usage()
		os.Exit(2)
	}

	if cfg.InputsFile != "" {
		if fileExists(cfg.DbFile) {
			fmt.Println("ERR: cannot have both inputs file and db file")
			os.Exit(2)
		}
	}

	/////////////////
	// init db

	db := initDb(cfg.DbFile)
	defer db.Close()

	if cfg.InputsFile != "" {
		// start fresh when inputs file is present
		createTables(db)
		urls := extractUrls(cfg.InputsFile)
		saveUrls(db, urls)
	}

	/////////////////
	// do tasks

	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}
	stats := new(Stats)

	tasks := loadTasks(db)
	bar := progressbar.New(len(tasks))
	sem := make(chan bool, cfg.Concurrency)
	for _, task := range tasks {
		sem <- true
		go func() {
			result := doRequest(task, httpClient, cfg.UserAgent)
			saveResult(db, result)
			updateStats(stats, result)
			bar.Describe(describeStats(stats))
			bar.Add(1)
			<-sem
		}()
	}
	for range cfg.Concurrency {
		sem <- true
	}
}
