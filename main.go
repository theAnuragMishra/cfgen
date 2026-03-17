package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli/v3"
)

type Problem struct {
	Code       string
	Statement  string
	InputSpec  string
	OutputSpec string
	InputEx    string
	OutputEx   []string
}

func main() {
	app := &cli.Command{
		Name:  "cfgen",
		Usage: "Fetch Codeforces problems and generate boilerplate",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "Author handle to include in boilerplate",
				Value:   "",
			},
			&cli.IntFlag{
				Name:    "workers",
				Aliases: []string{"w"},
				Usage:   "Number of concurrent workers",
				Value:   5,
			},
			&cli.BoolWithInverseFlag{
				Name:    "save-ps",
				Usage:   "Whether to create a problem statement file",
				Aliases: []string{"sps"},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			contestID := c.Args().Get(0)
			username := c.String("user")
			workers := c.Int("workers")
			ps := c.Bool("save-ps")
			return fetchContest(contestID, username, workers, ps)
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func fetchContest(contestID, user string, maxWorkers int, savePS bool) error {
	problemCodes, contestName, err := fetchProblemCodes(contestID)
	if err != nil {
		return err
	}

	fmt.Printf("📘 Contest: %s | Problems: %d | User: %s\n", contestName, len(problemCodes), user)
	os.MkdirAll(contestName, os.ModePerm)

	var wg sync.WaitGroup
	sema := make(chan struct{}, maxWorkers)

	for _, code := range problemCodes {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			sema <- struct{}{}
			defer func() { <-sema }()

			prob, err := getProblemDetails(contestID, code)
			if err != nil {
				log.Printf("⚠️ Error fetching problem %s: %v\n", code, err)
				return
			}
			err = saveProblem(contestName, prob, user, savePS)
			if err != nil {
				log.Printf("⚠️ Error saving problem %s: %v\n", code, err)
			} else {
				fmt.Printf("✅ Saved %s\n", code)
			}
		}(code)
	}
	wg.Wait()
	fmt.Println("🎉 All problems fetched and saved.")
	return nil
}

func fetchProblemCodes(contestID string) ([]string, string, error) {
	url := "https://codeforces.com/contest/" + contestID
	doc, err := fetchHTML(url)
	if err != nil {
		return nil, "", err
	}
	contestName := strings.TrimSpace(doc.Find(".rtable .left").First().Text())
	contestName = strings.ReplaceAll(contestName, ":", "")

	seen := make(map[string]bool)
	var codes []string

	doc.Find(`select[name="submittedProblemIndex"] option`).Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("value")
		if exists && len(value) != 0 && len(value) <= 2 && !seen[value] {
			seen[value] = true
			codes = append(codes, value)
		}
	})
	return codes, contestName, nil
}

func getProblemDetails(contestID, problemCode string) (Problem, error) {
	url := fmt.Sprintf("https://codeforces.com/contest/%s/problem/%s", contestID, problemCode)
	doc, err := fetchHTML(url)
	if err != nil {
		return Problem{}, err
	}
	statement := strings.TrimSpace(doc.Find(".problem-statement > div > p").Not(".input-specification p, .output-specification p, .sample-tests p, .note p").Text())
	inputSpec := strings.Join(doc.Find(".input-specification p").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	}), " ")
	outputSpec := strings.Join(doc.Find(".output-specification p").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	}), " ")
	inputEx := strings.Join(doc.Find(".sample-tests .input pre div").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	}), "\n")
	outputEx := strings.Split(strings.TrimSpace(doc.Find(".sample-tests .output pre").Text()), "\n")

	return Problem{
		Code:       problemCode,
		Statement:  statement,
		InputSpec:  inputSpec,
		OutputSpec: outputSpec,
		InputEx:    inputEx,
		OutputEx:   outputEx,
	}, nil
}

func saveProblem(contestName string, problem Problem, author string, savePS bool) error {
	path := fmt.Sprintf("%s/%s", contestName, problem.Code)
	os.MkdirAll(path, os.ModePerm)

	write := func(filename, content string) {
		err := os.WriteFile(fmt.Sprintf("%s/%s", path, filename), []byte(content), 0o644)
		if err != nil {
			log.Printf("Error writing %s: %v\n", filename, err)
		}
	}

	if savePS {
		statement := fmt.Sprintf("Problem Statement :-\n\n%s\n\n\nInput Specification :-\n\n%s\n\n\nOutput Specification :-\n\n%s\n\n\nInput Example :-\n\n%s\n\n\nOutput Example :-\n\n%s",
			problem.Statement, problem.InputSpec, problem.OutputSpec, problem.InputEx, strings.Join(problem.OutputEx, "\n"))
		write("problemStatement.txt", statement)
	}
	write("inputf.in", problem.InputEx)
	write("expectedf.out", strings.Join(problem.OutputEx, "\n"))
	write("outputf.out", "")
	write("solution.cpp", generateBoilerplate(author))
	return nil
}

func fetchHTML(url string) (*goquery.Document, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Spoof a browser and add common headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) "+
		"Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	var resp *http.Response
	// Retry loop for transient failures (rate limits, 5xx, etc.)
	for attempt := 0; attempt < 3; attempt++ {
		resp, err = client.Do(req)
		if err != nil {
			if attempt < 2 {
				time.Sleep(time.Duration(300*(1<<attempt)) * time.Millisecond)
				continue
			}
			return nil, err
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		// Potentially transient statuses — retry
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden || (resp.StatusCode >= 500 && resp.StatusCode < 600) {
			resp.Body.Close()
			if attempt < 2 {
				time.Sleep(time.Duration(300*(1<<attempt)) * time.Millisecond)
				continue
			}
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
		}

		// Non-retryable non-200 status: include a short body snippet for debugging
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s (%s)", resp.StatusCode, url, strings.TrimSpace(string(bodyBytes)))
	}

	if resp == nil {
		return nil, fmt.Errorf("no response received from %s", url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return doc, nil
}
