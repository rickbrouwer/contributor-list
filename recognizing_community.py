import subprocess
from collections import defaultdict

# Configuration variables
NUM_RELEASES = 20  # Number of releases to analyze
MIN_RELEASES = 3  # Minimum number of releases in which an author must have commits
NUM_COMMITS = 1  # Minimum number of commits per release

def get_releases(limit=NUM_RELEASES):
    """Fetch the latest 'limit' release tags from the repository, sorted by creation date."""
    result = subprocess.run(["git", "tag", "--list", "--sort=-creatordate"], capture_output=True, text=True)
    return result.stdout.splitlines()[:limit]

def get_commits_for_release(release, previous_release):
    """Fetch only commits that were added in this release compared to the previous release."""
    if previous_release:
        diff_range = f"{previous_release}..{release}"
    else:
        diff_range = release  # First release, count all commits
    
    result = subprocess.run(["git", "log", diff_range, "--pretty=format:%an"], capture_output=True, text=True)
    return result.stdout.splitlines()

def count_commits_per_author(limit=NUM_RELEASES):
    """Count commits per GitHub author per release, limited to the last 'limit' releases."""
    releases = get_releases(limit)
    author_contributions = defaultdict(lambda: defaultdict(int))
    
    for i, release in enumerate(releases):
        previous_release = releases[i + 1] if i + 1 < len(releases) else None
        authors = get_commits_for_release(release, previous_release)
        for author in authors:
            author_contributions[author][release] += 1
    
    return author_contributions

def filter_active_authors(author_contributions, min_releases=MIN_RELEASES, num_commits=NUM_COMMITS):
    """Filter authors who have at least 'num_commits' in 'min_releases' different releases."""
    return {author: releases for author, releases in author_contributions.items() if sum(1 for r in releases if releases[r] >= num_commits) >= min_releases}

def main():
    author_contributions = count_commits_per_author()
    active_authors = filter_active_authors(author_contributions)
    
    print(f"Contributors with at least {NUM_COMMITS} commits in at least {MIN_RELEASES} of the last {NUM_RELEASES} releases:")
    for author, releases in active_authors.items():
        release_info = ', '.join(f"{r} ({releases[r]} commits)" for r in releases)
        print(f"{author}: {release_info}")

if __name__ == "__main__":
    main()
