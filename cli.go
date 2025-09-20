package main

import (
	"flag"
	"fmt"
)

func parseArgs() *Config {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `fetch2 downloads web pages from URLs and saves them to a database.

Usage:
  fetch2 [flags]

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), `
Notes:
  -inputs cannot be used if db file already exists. When used, it creates the database and loads URLs from the input file.
  -inputs can have URLs anywhere in each line. The tool finds all http and https links.

Examples:
  fetch2 -db a.db -inputs urls.txt
  fetch2 -db a.db -concurrency 10 -timeout 10 -user-agent 'MyAgent/1.0'
`)
	}

	cfg := &Config{
		UserAgent:   *flag.String("user-agent", "fetch2", "Name sent to websites when making requests"),
		InputsFile:  *flag.String("inputs", "", "File with URLs to fetch. Starts fresh database. Cannot use if db file exists."),
		DbFile:      *flag.String("db", "", "Database file path (required). Creates if not exists."),
		Concurrency: *flag.Int("concurrency", 6, "How many downloads to run at once"),
		Timeout:     *flag.Int("timeout", 30, "How long to wait for each download (seconds)"),
	}

	flag.Parse()
	return cfg
}
