package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Configuration variables
const (
	NUM_RELEASES           = 20    // Number of releases to analyze
	MIN_RELEASES           = 3     // Minimum number of releases in which an author must have commits
	NUM_COMMITS            = 1     // Minimum number of commits per release
	SHOW_RECENT_QUALIFIERS = false // Show only authors who just met the MIN_RELEASES criteria in the latest release
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
func countCommitsPerAuthor(limit int) (map[string]map[string]int, []string) {
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

	return authorContributions, releases
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

// identifyRecentQualifiers identifies authors who exactly met the minimum release criteria with the latest release
func identifyRecentQualifiers(authorContributions map[string]map[string]int, releases []string, minReleases, numCommits int) map[string]map[string]int {
	recentQualifiers := make(map[string]map[string]int)

	if len(releases) == 0 {
		return recentQualifiers
	}

	latestRelease := releases[0]

	for author, releaseData := range authorContributions {
		// Count releases where this author has sufficient commits
		releaseCount := 0
		hasLatestRelease := false

		for release, commitCount := range releaseData {
			if commitCount >= numCommits {
				releaseCount++
				if release == latestRelease {
					hasLatestRelease = true
				}
			}
		}

		// Check if author exactly met the threshold with the latest release
		if releaseCount >= minReleases && hasLatestRelease {
			// Remove the latest release and check if the author would then fail to meet criteria
			testCount := 0
			for release, commitCount := range releaseData {
				if release != latestRelease && commitCount >= numCommits {
					testCount++
				}
			}

			if testCount < minReleases {
				recentQualifiers[author] = releaseData
			}
		}
	}

	return recentQualifiers
}
func main() {
    authorContributions, releases := countCommitsPerAuthor(NUM_RELEASES)
    
    if SHOW_RECENT_QUALIFIERS {
        recentQualifiers := identifyRecentQualifiers(authorContributions, releases, MIN_RELEASES, NUM_COMMITS)
        
        if len(recentQualifiers) == 0 {
            fmt.Printf("No new contributors who just met the criteria in the latest release (%d commits in %d releases)\n", 
                NUM_COMMITS, MIN_RELEASES)
        } else {
            fmt.Printf("New contributors who just met the criteria in the latest release (%d commits in %d releases):\n", 
                NUM_COMMITS, MIN_RELEASES)
            
            for author, releaseData := range recentQualifiers {
                var releaseInfo []string
                for release, count := range releaseData {
                    releaseInfo = append(releaseInfo, fmt.Sprintf("%s (%d commits)", release, count))
                }
                fmt.Printf("%s: %s\n", author, strings.Join(releaseInfo, ", "))
            }
        }
    } else {
        activeAuthors := filterActiveAuthors(authorContributions, MIN_RELEASES, NUM_COMMITS)
        
        if len(activeAuthors) == 0 {
            fmt.Printf("No contributors with at least %d commits in at least %d of the last %d releases\n", 
                NUM_COMMITS, MIN_RELEASES, NUM_RELEASES)
        } else {
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
    }
}
