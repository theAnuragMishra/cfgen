package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
		Commands: []*cli.Command{
			{
				Name:  "fetch",
				Usage: "Download problems from a Codeforces contest",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Aliases:  []string{"c"},
						Usage:    "Codeforces contest ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "user",
						Aliases:  []string{"u"},
						Usage:    "Author handle to include in boilerplate",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "workers",
						Aliases: []string{"w"},
						Usage:   "Number of concurrent workers",
						Value:   5,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					contestID := strconv.Itoa(c.Int("id"))
					username := c.String("user")
					workers := c.Int("workers")

					return fetchContest(contestID, username, workers)
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func fetchContest(contestID, user string, maxWorkers int) error {
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
			err = saveProblem(contestName, prob, user)
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

func saveProblem(contestName string, problem Problem, author string) error {
	path := fmt.Sprintf("%s/%s", contestName, problem.Code)
	os.MkdirAll(path, os.ModePerm)

	write := func(filename, content string) {
		err := os.WriteFile(fmt.Sprintf("%s/%s", path, filename), []byte(content), 0o644)
		if err != nil {
			log.Printf("Error writing %s: %v\n", filename, err)
		}
	}

	statement := fmt.Sprintf("Problem Statement :-\n\n%s\n\n\nInput Specification :-\n\n%s\n\n\nOutput Specification :-\n\n%s\n\n\nInput Example :-\n\n%s\n\n\nOutput Example :-\n\n%s",
		problem.Statement, problem.InputSpec, problem.OutputSpec, problem.InputEx, strings.Join(problem.OutputEx, "\n"))

	write("problemStatement.txt", statement)
	write("inputf.in", problem.InputEx)
	write("expectedf.out", strings.Join(problem.OutputEx, "\n"))
	write("outputf.out", "")
	write("solution.cpp", generateBoilerplate(author))
	return nil
}

func generateBoilerplate(name string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(`/**
 *    author: %s
 *    created: %s
**/
#include <bits/stdc++.h>

#define for0(i, n) for (int i = 0; i < (int)(n); ++i)
#define for1(i, n) for (int i = 1; i <= (int)(n); ++i)
#define forc(i, l, r) for (int i = (int)(l); i <= (int)(r); ++i)
#define forr0(i, n) for (int i = (int)(n) - 1; i >= 0; --i)
#define forr1(i, n) for (int i = (int)(n); i >= 1; --i)
#define each(x, a) for (auto &x : a)

#define pb push_back
#define fi first
#define se second
#define eb emplace_back
#define ef emplace_front
#define em emplace
#define fr front()
#define bk back()

#define bpc __builtin_popcount
#define bpcll __builtin_popcountll
#define clz __builtin_clz
#define clzll __builtin_clzll
#define ctzll __builtin_ctzll
#define ctz __builtin_ctz
#define sqrt __builtin_sqrt
#define abs __builtin_abs
#define memset __builtin_memset
#define memcpy __builtin_memcpy

#define all(x) (x).begin(), (x).end()
#define rall(x) (x).rbegin(), (x).rend()

#define present(c, x) ((c).find(x) != (c).end())
#define cpresent(c, x) (find(all(c), x) != c.end())

#define wne(c) while (!((c).empty()))

#define sz(a) int((a).size())

using namespace std;

using ll = long long;
using db = double;
using ld = long double;
using ul = unsigned long;
using ull = unsigned long long;

using vi = vector<int>;
using vvi = vector<vi>;
using pi = pair<int, int>;
using pll = pair<ll, ll>;
using pdb = pair<db, db>;
using vpi = vector<pi>;
using vc = vector<char>;
using vdb = vector<db>;
using vs = vector<string>;
using vll = vector<ll>;
using vvll = vector<vll>;
using vb = vector<bool>;
using si = unordered_set<int>;
using mi = unordered_map<int, int>;
using sc = unordered_set<char>;
using pqi = priority_queue<int>;
using pqpi = priority_queue<pi>;
using pqll = priority_queue<ll>;
using pqpll = priority_queue<pll>;

#define MOD 1000000007

#define sp << ' ' <<
#define br << '\n'

void solve() {}

int main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
#ifndef ONLINE_JUDGE
  freopen("inputf.in", "r", stdin);
  freopen("outputf.out", "w", stdout);
#endif
  int t = 1;
  cin >> t;
  while (t--) {
    solve();
  }
}`, name, now)
}

func fetchHTML(url string) (*goquery.Document, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Spoof a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) "+
		"Chrome/115.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return doc, nil
}
