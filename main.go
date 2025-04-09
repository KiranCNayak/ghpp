package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/fatih/color"
)

// ---- Struct Definitions ----

type Owner struct {
	Login string `json:"login"`
}

type License struct {
	Name string `json:"name"`
}

type Repo struct {
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	HTMLURL         string    `json:"html_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	StargazersCount int       `json:"stargazers_count"`
	Forks           int       `json:"forks"`
	Watchers        int       `json:"watchers"`
	Size            int       `json:"size"`
	Owner           Owner     `json:"owner"`
	License         *License  `json:"license"`
}

// ---- Configurable Defaults ----

var defaultFields = []string{
	"name",
	"full_name",
	"html_url",
	"created_at",
	"updated_at",
	"stargazers_count",
}

// ---- Main ----

func main() {
	include := flag.String("include", "", "Comma-separated fields to include")
	exclude := flag.String("exclude", "", "Comma-separated fields to exclude")
	since := flag.Bool("since", false, "Show created_at as 'X years Y months Z days ago'")
	short := flag.Bool("short", false, "Show time difference in short format (e.g., 2y 5m ago)")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: ghpp <owner>/<repo> [--include=\"\"] [--exclude=\"\"] [--since]")
		os.Exit(1)
	}

	repoArg := flag.Arg(0)
	parts := strings.Split(repoArg, "/")
	if len(parts) != 2 {
		fmt.Println("Invalid repo format. Use <owner>/<repo>")
		os.Exit(1)
	}

	owner, repo := parts[0], parts[1]
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Failed to fetch repository:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("GitHub API returned status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var repoData Repo
	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		fmt.Println("Failed to parse JSON:", err)
		os.Exit(1)
	}

	fieldsToShow := buildFieldList(defaultFields, *include, *exclude)
	printRepoInfo(repoData, fieldsToShow, *since || *short, *short)
}

// ---- Field Logic ----

func buildFieldList(defaults []string, include, exclude string) []string {
	fieldSet := make(map[string]bool)

	for _, f := range defaults {
		fieldSet[f] = true
	}
	if include != "" {
		for _, f := range strings.Split(include, ",") {
			fieldSet[strings.TrimSpace(f)] = true
		}
	}
	if exclude != "" {
		for _, f := range strings.Split(exclude, ",") {
			delete(fieldSet, strings.TrimSpace(f))
		}
	}

	result := []string{}
	for f := range fieldSet {
		result = append(result, f)
	}
	return result
}

// ---- Pretty Printing ----

func printRepoInfo(repo Repo, fields []string, since bool, short bool) {
	for _, field := range fields {
		switch field {
		case "name":
			color.Cyan("📦 Name: %s", repo.Name)
		case "full_name":
			color.Cyan("📛 Full Name: %s", repo.FullName)
		case "html_url":
			color.Blue("🌐 URL: %s", repo.HTMLURL)
		case "created_at":
			if since {
				if short {
					color.Green("📅 Created: %s", humanTimeDiffShort(repo.CreatedAt))
				} else {
					color.Green("📅 Created: %s", humanTimeDiff(repo.CreatedAt))
				}
			} else {
				color.Green("📅 Created: %s", repo.CreatedAt.Format(time.RFC3339))
			}
		case "updated_at":
			color.Yellow("🔄 Updated: %s", repo.UpdatedAt.Format(time.RFC3339))
		case "stargazers_count":
			color.Magenta("⭐ Stars: %d", repo.StargazersCount)
		case "forks":
			color.Magenta("🍴 Forks: %d", repo.Forks)
		case "watchers":
			color.Magenta("👀 Watchers: %d", repo.Watchers)
		case "size":
			color.Magenta("📦 Size: %d KB", repo.Size)
		case "owner.login":
			color.Cyan("👤 Owner: %s", repo.Owner.Login)
		case "license.name":
			if repo.License != nil {
				color.Cyan("📝 License: %s", repo.License.Name)
			}
		default:
			color.Red("❓ Unknown field: %s", field)
		}
	}
}

// ---- Time Difference (Since) ----
func humanTimeDiff(t time.Time) string {
	now := time.Now()
	years := now.Year() - t.Year()
	months := int(now.Month()) - int(t.Month())
	days := now.Day() - t.Day()

	if days < 0 {
		months--
		days += 30 // Approximate
	}
	if months < 0 {
		years--
		months += 12
	}

	parts := []string{}
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d year%s", years, plural(years)))
	}
	if months > 0 {
		parts = append(parts, fmt.Sprintf("%d month%s", months, plural(months)))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d day%s", days, plural(days)))
	}

	if len(parts) == 0 {
		return "today"
	}

	return strings.Join(parts, " ") + " ago"
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// ---- Time Difference (Short) ----
func humanTimeDiffShort(t time.Time) string {
	now := time.Now()
	years := now.Year() - t.Year()
	months := int(now.Month()) - int(t.Month())
	days := now.Day() - t.Day()

	if days < 0 {
		months--
		days += 30
	}
	if months < 0 {
		years--
		months += 12
	}

	parts := []string{}
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%dy", years))
	}
	if months > 0 {
		parts = append(parts, fmt.Sprintf("%dm", months))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}

	if len(parts) == 0 {
		return "today"
	}

	return strings.Join(parts, " ") + " ago"
}
