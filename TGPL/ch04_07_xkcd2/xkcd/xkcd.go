package xkcd

import "html/template"

// IssuesURL ...
const xkcdURL = "https://xkcd.com/%d/info.0.json"

// ComicInfo ...
type ComicInfo struct {
	Num        int
	Title      string
	SafeTitle  string `json:"safe_title"`
	Transcript string
	Alt        string
	Day        string
	Month      string
	Year       string
	Img        template.URL
	News       string
}
