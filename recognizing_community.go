package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Configuration variables
const (
	NUM_RELEASES = 20 // Number of releases to analyze
	MIN_RELEASES = 3  // Minimum number of releases in which an author must have commits
	NUM_COMMITS  = 1  // Minimum number of commits per release
)

// getReleases fetches the latest 'limit' release tags from the repository, sorted by creation date
func getReleases(limit int) []string {
	cmd := exec.Command("git", "tag", "--list", "--sort=-creatordate")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error getting releases: %v\n", err)
		return []string{}
	}

	releases := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(releases) > limit {
		releases = releases[:limit]
	}
	return releases
}

// getCommitsForRelease fetches only commits that were added in this release compared to the previous release
func getCommitsForRelease(release, previousRelease string) []string {
	var diffRange string
	if previousRelease != "" {
		diffRange = previousRelease + ".." + release
	} else {
		diffRange = release // First release, count all commits
	}

	cmd := exec.Command("git", "log", diffRange, "--pretty=format:%an")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error getting commits for release %s: %v\n", release, err)
		return []string{}
	}

	return strings.Split(strings.TrimSpace(out.String()), "\n")
}

// countCommitsPerAuthor counts commits per GitHub author per release, limited to the last 'limit' releases
func countCommitsPerAuthor(limit int) map[string]map[string]int {
	releases := getReleases(limit)
	authorContributions := make(map[string]map[string]int)

	for i, release := range releases {
		var previousRelease string
		if i+1 < len(releases) {
			previousRelease = releases[i+1]
		}

		authors := getCommitsForRelease(release, previousRelease)
		for _, author := range authors {
			if author == "" {
				continue
			}
			
			if _, exists := authorContributions[author]; !exists {
				authorContributions[author] = make(map[string]int)
			}
			authorContributions[author][release]++
		}
	}

	return authorContributions
}

// filterActiveAuthors filters authors who have at least 'numCommits' in 'minReleases' different releases
func filterActiveAuthors(authorContributions map[string]map[string]int, minReleases, numCommits int) map[string]map[string]int {
	activeAuthors := make(map[string]map[string]int)

	for author, releases := range authorContributions {
		releaseCount := 0
		for _, count := range releases {
			if count >= numCommits {
				releaseCount++
			}
		}

		if releaseCount >= minReleases {
			activeAuthors[author] = releases
		}
	}

	return activeAuthors
}

func main() {
	authorContributions := countCommitsPerAuthor(NUM_RELEASES)
	activeAuthors := filterActiveAuthors(authorContributions, MIN_RELEASES, NUM_COMMITS)

	fmt.Printf("Contributors with at least %d commits in at least %d of the last %d releases:\n", 
		NUM_COMMITS, MIN_RELEASES, NUM_RELEASES)
	
	for author, releases := range activeAuthors {
		var releaseInfo []string
		for release, count := range releases {
			releaseInfo = append(releaseInfo, fmt.Sprintf("%s (%d commits)", release, count))
		}
		fmt.Printf("%s: %s\n", author, strings.Join(releaseInfo, ", "))
	}
}
