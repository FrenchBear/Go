package xkcd

// IssuesURL ...
const xkcdURL = "https://xkcd.com/%d/info.0.json"

// ComicInfo ...
type ComicInfo struct {
	Title      string
	SafeTitle  string `json:"safe_title"`
	Transcript string
	Alt        string
}
