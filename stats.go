package main

import "fmt"

func updateStats(stats *Stats, result *Result) {
	if result.Err != nil {
		stats.Err.Add(1)
		return
	}
	status := *result.Status
	if status >= 200 && status < 300 {
		stats.Status200.Add(1)
	} else if status >= 400 && status < 500 {
		stats.Status400.Add(1)
	} else if status >= 500 && status < 600 {
		stats.Status500.Add(1)
	}
}

func describeStats(stats *Stats) string {
	return fmt.Sprintf("[200] %d, [400] %d, [500] %d, [err] %d", stats.Status200.Load(), stats.Status400.Load(), stats.Status500.Load(), stats.Err.Load())
}
