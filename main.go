package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

const upStatus = "UP"
const downStatus = "DOWN"

type CheckResult struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Success    bool
}

func checkURL(client *http.Client, url string) (CheckResult, error) {
	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	result := CheckResult{
		URL:      url,
		Duration: duration,
	}

	if err != nil {
		return result, err
	}

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	defer resp.Body.Close()

	return result, nil
}

func main() {

	interval := flag.Duration(
		"interval",
		2*time.Second,
		"Интервал между проверками",
	)

	timeout := flag.Duration(
		"timeout",
		10*time.Second,
		"Максимальное время HTTP-запроса",
	)
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Use: go run . [flags] <url>")
		flag.PrintDefaults()
		return
	}

	url := flag.Arg(0)

	client := &http.Client{
		Timeout: *timeout,
	}

	for {
		result, err := checkURL(client, url)
		checkedAt := time.Now().Format(time.RFC3339)
		if !result.Success {
			if err != nil {
				fmt.Printf("%s %s Request failed: %v\n", checkedAt, downStatus, err)
			} else {
				fmt.Println(checkedAt, downStatus, result.StatusCode, result.Duration.Round(time.Millisecond), result.URL)
			}

		} else {
			fmt.Println(checkedAt, upStatus, result.StatusCode, result.Duration.Round(time.Millisecond), result.URL)
		}
		time.Sleep(*interval)
	}

}
