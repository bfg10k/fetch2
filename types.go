package main

import "sync/atomic"

type Config struct {
	UserAgent   string
	InputsFile  string
	DbFile      string
	Concurrency int
	Timeout     int
}

type Task struct {
	Id  int
	Url string
}

type Result struct {
	Id      int
	Status  *int
	Headers *string
	Content []byte
	Err     *string
}

type Stats struct {
	Status200 atomic.Int64
	Status400 atomic.Int64
	Status500 atomic.Int64
	Err       atomic.Int64
}
