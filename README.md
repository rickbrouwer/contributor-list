# Git Release Contribution Analyzer

Open-source project are only as successful as the community that drives it. To show your appreciation to essential members of your community, you can recognize them by issuing Credly badges.  
This Python script analyzes contributions to a Git repository by counting the number of commits per author in the last few releases. It identifies contributors who meet a specified threshold of commits across multiple releases.

## Features
- Fetches the latest release tags.
- Counts commits per author for each release.
- Filters authors who meet the required commit threshold.
- Outputs contributors who have been active across multiple releases.

## Configuration
You can modify the following variables at the top of the script to adjust the analysis:
- `NUM_RELEASES`: Number of latest releases to analyze.
- `MIN_RELEASES`: Minimum number of releases an author must have contributed to.
- `NUM_COMMITS`: Minimum number of commits required per release.

## Usage
Run the script in a cloned Git repository:
```sh
python script.py
```

## Example Output
```
Contributors with at least 1 commits in at least 3 of the last 20 releases:
username1: v1.2 (3 commits), v1.3 (5 commits), v1.4 (2 commits)
username2: v1.1 (4 commits), v1.2 (2 commits), v1.3 (6 commits)
```

## Requirements
- Python 3
- Git installed and accessible via the command line

## Notes
Ensure you run the script **inside** a valid Git repository containing tags for the releases you want to analyze.

