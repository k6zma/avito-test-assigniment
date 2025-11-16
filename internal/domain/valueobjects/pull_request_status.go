package valueobjects

type PullRequestStatus string

const (
	OpenPullRequest   PullRequestStatus = "OPEN"
	MergedPullRequest PullRequestStatus = "MERGED"
)
